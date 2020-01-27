package client

import (
	"context"
	"github.com/ducc/profile-collector/protos"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

// default option values
var defaultBuilder = builder{
	address:    "", // todo default address
	maxRetries: 3,
	backoff:    time.Second * 1,
	useTracing: true,
}

func New(ctx context.Context, opts ...Option) (protos.StoreClient, error) {
	builder := defaultBuilder
	for _, opt := range opts {
		opt(&builder)
	}

	logrus.WithField("address", builder.address).Debug("connecting to store server grpc")

	retryOptions := []grpc_retry.CallOption{
		grpc_retry.WithMax(builder.maxRetries),
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(builder.backoff)),
	}

	unaryInterceptors := []grpc.UnaryClientInterceptor{
		grpc_prometheus.UnaryClientInterceptor,
		grpc_retry.UnaryClientInterceptor(retryOptions...),
	}

	streamInterceptors := []grpc.StreamClientInterceptor{
		grpc_prometheus.StreamClientInterceptor,
		grpc_retry.StreamClientInterceptor(retryOptions...),
	}

	if builder.useTracing && opentracing.IsGlobalTracerRegistered() {
		unaryInterceptors = append(unaryInterceptors, otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()))
		streamInterceptors = append(streamInterceptors, otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer()))
	}

	conn, err := grpc.DialContext(
		ctx, builder.address,
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(unaryInterceptors...),
		grpc.WithChainStreamInterceptor(streamInterceptors...),
	)
	if err != nil {
		return nil, err
	}

	return protos.NewStoreClient(conn), nil
}

type Option func(*builder)

type builder struct {
	address    string
	maxRetries uint
	backoff    time.Duration
	useTracing bool
}
