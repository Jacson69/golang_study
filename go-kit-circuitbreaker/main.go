package main

import (
	"context"
	"flag"
	"fmt"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/juju/ratelimit"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang_study/go-kit-circuitbreaker/endpoint"
	"golang_study/go-kit-circuitbreaker/instrument"
	"golang_study/go-kit-circuitbreaker/plugins"
	"golang_study/go-kit-circuitbreaker/register"
	"golang_study/go-kit-circuitbreaker/service"
	"golang_study/go-kit-circuitbreaker/transport"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var (
		consulHost  = flag.String("consul.host", "192.168.78.172", "consul host")
		consulPort  = flag.String("consul.port", "8500", "consul port")
		serviceHost = flag.String("service.host", "192.168.0.20", "service host")
		servicePort = flag.String("service.port", "9000", "service port")
		zipkinURL   = flag.String("zipkin.url", "http://192.168.78.172:9411/api/v2/spans", "Zipkin server url")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)

		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var svc service.Service
	svc = service.ArithmeticService{}

	//{
	// add logging middleware
	// 开启日志功能
	svc = plugins.LoggingMiddleware(logger)(svc)

	//链路追踪
	var zipkinTracer *zipkin.Tracer
	{
		var (
			err           error
			hostPort      = *serviceHost + ":" + *servicePort
			serviceName   = "arithmetic-service"
			useNoopTracer = (*zipkinURL == "")
			reporter      = zipkinhttp.NewReporter(*zipkinURL)
		)
		defer reporter.Close()
		zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
		zipkinTracer, err = zipkin.NewTracer(
			reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
		)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		if !useNoopTracer {
			logger.Log("tracer", "Zipkin", "type", "Native", "URL", *zipkinURL)
		}
	}

	// 开启Prometheus监控功能
	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "raysonxin",
		Subsystem: "arithmetic_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "raysonxin",
		Subsystem: "arithemetic_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	svc = instrument.Metrics(requestCount, requestLatency)(svc)

	//{
	//artithmeticEndpoint := endpoint.MakeArithmeticEndpoint(svc)
	//
	//// add ratelimit,refill every second,set capacity 3
	//// 在此处添加限流的功能
	//ratebucket := ratelimit.NewBucket(time.Second*1, 100)
	//artithmeticEndpoint = instrument.NewTokenBucketLimitterWithJuju(ratebucket)(artithmeticEndpoint)
	////}
	//
	////添加追踪，设置span的名称为calculate-endpoint
	//artithmeticEndpoint = kitzipkin.TraceEndpoint(zipkinTracer, "calculate-endpoint")(artithmeticEndpoint)
	//}

	calEndpoint := endpoint.MakeArithmeticEndpoint(svc)
	// add ratelimit,refill every second,set capacity 3
	// 在此处添加限流的功能
	ratebucket := ratelimit.NewBucket(time.Second*1, 100)
	calEndpoint = instrument.NewTokenBucketLimitterWithJuju(ratebucket)(calEndpoint)
	//添加追踪，设置span的名称为calculate-endpoint
	calEndpoint = kitzipkin.TraceEndpoint(zipkinTracer, "calculate-endpoint")(calEndpoint)
	calEndpoint = kitjwt.NewParser(service.JwtKeyFunc, jwt2.SigningMethodHS256, kitjwt.StandardClaimsFactory)(calEndpoint)

	//创建健康检查的Endpoint，增加了限流
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)
	healthEndpoint = instrument.NewTokenBucketLimitterWithJuju(ratebucket)(healthEndpoint)
	healthEndpoint = kitzipkin.TraceEndpoint(zipkinTracer, "health-endpoint")(healthEndpoint)

	//身份认证Endpoint
	authEndpoint := endpoint.MakeAuthEndpoint(svc)
	authEndpoint = instrument.NewTokenBucketLimitterWithJuju(ratebucket)(authEndpoint)
	authEndpoint = kitzipkin.TraceEndpoint(zipkinTracer, "login-endpoint")(authEndpoint)

	//把算术运算Endpoint和健康检查Endpoint和身份认证AuthEndpoint封装至ArithmeticEndpoints
	endpts := endpoint.ArithmeticEndpoints{
		//ArithmeticEndpoint:  artithmeticEndpoint,
		ArithmeticEndpoint:  calEndpoint,
		HealthCheckEndpoint: healthEndpoint,
		AuthEndpoint:        authEndpoint,
	}

	r := transport.MakeHttpHandler(ctx, endpts, zipkinTracer, logger)

	registar := register.Register(*consulHost, *consulPort, *serviceHost, *servicePort, logger)

	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		//启动前执行注册
		registar.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出，取消注册
	registar.Deregister()
	fmt.Println(error)
}
