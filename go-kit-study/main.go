package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/juju/ratelimit"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang_study/go-kit-study/endpoint"
	"golang_study/go-kit-study/instrument"
	"golang_study/go-kit-study/plugins"
	"golang_study/go-kit-study/register"
	"golang_study/go-kit-study/service"
	"golang_study/go-kit-study/transport"
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

	artithmeticEndpoint := endpoint.MakeArithmeticEndpoint(svc)
	// add ratelimit,refill every second,set capacity 3
	// 在此处添加限流的功能
	ratebucket := ratelimit.NewBucket(time.Second*1, 100)
	artithmeticEndpoint = instrument.NewTokenBucketLimitterWithJuju(ratebucket)(artithmeticEndpoint)
	//}

	//创建健康检查的Endpoint，未增加限流
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	//把算术运算Endpoint和健康检查Endpoint封装至ArithmeticEndpoints
	endpts := endpoint.ArithmeticEndpoints{
		ArithmeticEndpoint:  artithmeticEndpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	r := transport.MakeHttpHandler(ctx, endpts, logger)

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
