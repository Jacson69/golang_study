package instrument

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/juju/ratelimit"
	"golang_study/go-kit-study/service"
	"time"
)

// 该文件主要是实现限流的功能
var ErrLimitExceed = errors.New("Rate limit exceed!")

// NewTokenBucketLimitterWithJuju 使用juju/ratelimit创建限流中间件
func NewTokenBucketLimitterWithJuju(bkt *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if bkt.TakeAvailable(1) == 0 {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

// 以下是监控功能加入的代码
// metricMiddleware 定义监控中间件，嵌入Service
// 新增监控指标项：requestCount和requestLatency
type metricMiddleware struct {
	service.Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

// Metrics 指标采集方法
func Metrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return metricMiddleware{
			next,
			requestCount,
			requestLatency}
	}
}

func (mw metricMiddleware) Add(a, b int) (ret int) {

	defer func(beign time.Time) {
		lvs := []string{"method", "Add"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())

	ret = mw.Service.Add(a, b)
	return ret
}

func (mw metricMiddleware) Substract(a, b int) (ret int) {

	defer func(beign time.Time) {
		lvs := []string{"method", "Substract"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())

	ret = mw.Service.Substract(a, b)
	return ret
}

func (mw metricMiddleware) Multiply(a, b int) (ret int) {

	defer func(beign time.Time) {
		lvs := []string{"method", "Multiply"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())

	ret = mw.Service.Multiply(a, b)
	return ret
}

func (mw metricMiddleware) Divide(a, b int) (ret int, err error) {

	defer func(beign time.Time) {
		lvs := []string{"method", "Divide"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())

	ret, err = mw.Service.Divide(a, b)
	return ret, err
}

// metricMiddleware实现HealthCheck
func (mw metricMiddleware) HealthCheck() (result bool) {

	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result = mw.Service.HealthCheck()
	return
}
