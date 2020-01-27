package main

import (
	"context"
	"flag"
	"github.com/ducc/profile-collector/metrics"
	"github.com/ducc/profile-collector/protos"
	"github.com/ducc/profile-collector/store"
	"github.com/ducc/profile-collector/store/aggregator"
	"github.com/ducc/profile-collector/store/profiles"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

var (
	address   string
	logLevel  string
	gcsBucket string
)

func init() {
	flag.StringVar(&address, "address", ":8080", "grpc server address")
	flag.StringVar(&logLevel, "level", "debug", "logrus logging level")
	flag.StringVar(&gcsBucket, "gcsBucket", "", "gcs profile store bucket")
}

func main() {
	logging()
	go metrics.Start()

	ctx := context.Background()

	// profilesStore, err := profiles.NewGCS(ctx, gcsBucket)
	// if err != nil {
	// 	logrus.WithError(err).Fatal("creating profiles store")
	// }
	profilesStore := &profiles.MockProfiles{}

	aggStore, err := aggregator.NewPrometheus()
	if err != nil {
		logrus.WithError(err).Fatal("creating aggregator store")
	}

	server, err := store.New(profilesStore, aggStore)
	if err != nil {
		logrus.WithError(err).Fatal("creating store server")
	}

	if err := serve(server); err != nil {
		logrus.WithError(err).Fatal("serving requests")
	}
}

func serve(protoServer protos.StoreServer) error {
	logrus.WithField("address", address).Debug("running grpc server")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer()),
		)),
	)

	protos.RegisterStoreServer(grpcServer, protoServer)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			grpcServer.GracefulStop()
			os.Exit(1)
		}
	}()

	return grpcServer.Serve(listener)
}

func logging() {
	if ll, err := logrus.ParseLevel(logLevel); err != nil {
		logrus.WithError(err).Fatal("unable to parse logrus logging level")
	} else {
		logrus.SetLevel(ll)
	}
}
