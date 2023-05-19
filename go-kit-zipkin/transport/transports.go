package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang_study/go-kit-zipkin/endpoint"
	"net/http"
	"strconv"
)

// decodeArithmeticRequest decode request params to struct
func decodeArithmeticRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	requestType, ok := vars["type"]
	if !ok {
		return nil, ErrorBadRequest
	}

	pa, ok := vars["a"]
	if !ok {
		return nil, ErrorBadRequest
	}

	pb, ok := vars["b"]
	if !ok {
		return nil, ErrorBadRequest
	}

	a, _ := strconv.Atoi(pa)
	b, _ := strconv.Atoi(pb)

	return endpoint.ArithmeticRequest{
		RequestType: requestType,
		A:           a,
		B:           b,
	}, nil
}

// encodeArithmeticResponse encode response to return
func encodeArithmeticResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpoint.HealthRequest{}, nil
}

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

//当只有一个endpoint的时候用
// MakeHttpHandler make http handler use mux
//func MakeHttpHandler(ctx context.Context, endpoint kitendpoint.Endpoint, logger log.Logger) http.Handler {
//	r := mux.NewRouter()
//
//	options := []kithttp.ServerOption{
//		kithttp.ServerErrorLogger(logger),
//		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
//	}
//
//	r.Methods("POST").Path("/calculate/{type}/{a}/{b}").Handler(kithttp.NewServer(
//		endpoint,
//		decodeArithmeticRequest,
//		encodeArithmeticResponse,
//		options...,
//	))
//	r.Path("/metrics").Handler(promhttp.Handler())
//
//	return r
//}

// 以下是多个endpoint的时候用
func MakeHttpHandler(ctx context.Context, endpoints endpoint.ArithmeticEndpoints, zipkinTracer *gozipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	zipkinServer := zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name("http-transport"))
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		zipkinServer,
	}

	r.Methods("POST").Path("/calculate/{type}/{a}/{b}").Handler(kithttp.NewServer(
		endpoints.ArithmeticEndpoint,
		decodeArithmeticRequest,
		encodeArithmeticResponse,
		options...,
	))
	r.Path("/metrics").Handler(promhttp.Handler())

	// create health check handler
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeArithmeticResponse,
		options...,
	))

	return r
}
