package client

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zhangx1n/plato/common/prpc/prome"
	"github.com/zhangx1n/plato/common/prpc/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const nameSpace = "prpc_client"

var (
	clientHandleCounter = prome.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: nameSpace,
			Subsystem: "req",
			Name:      "client_handle_total",
		},
		[]string{"method", "server", "code", "ip"},
	)

	clientHandleHistogram = prome.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: nameSpace,
			Subsystem: "req",
			Name:      "client_handle_seconds",
		},
		[]string{"method", "server", "ip"},
	)
)

// MetricUnaryClientInterceptor ...
func MetricUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		code := status.Code(err)
		clientHandleCounter.WithLabelValues(method, cc.Target(), code.String(), util.ExternalIP()).Inc()
		clientHandleHistogram.WithLabelValues(method, cc.Target(), util.ExternalIP()).Observe(time.Since(beg).Seconds())

		return err
	}
}
