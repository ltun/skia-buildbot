// package main is the main entry point for the cabe server executable.
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"go.skia.org/infra/go/httputils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"go.skia.org/infra/go/common"
	"go.skia.org/infra/go/sklog"

	"go.skia.org/infra/cabe/go/analysisserver"
	"go.skia.org/infra/cabe/go/backends"
	cpb "go.skia.org/infra/cabe/go/proto"
)

const (
	appName = "cabe"
)

var (
	host     = flag.String("host", "localhost", "HTTP service host")
	port     = flag.Int("port", 8002, "http service port (e.g., '8002')")
	grpcPort = flag.Int("grpc_port", 50051, "gRPC service port (e.g., '50051')")
	promPort = flag.String("prom_port", ":20000", "Metrics service address (e.g., ':10110')")
)

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := authorize(ctx); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		sklog.Infof("incoming grpc metadata: %#v", md)
	} else {
		sklog.Infof("no incoming grpc metadata")
	}
	return nil
}

func main() {
	ctx := context.Background()
	flag.Parse()

	// Setup flags.
	common.InitWithMust(
		appName,
		common.PrometheusOpt(promPort),
	)

	go func() {
		// Just testing the http healthz check to make sure envoy can
		// connect to these processes at all. If we end up needing
		// both the http server and the grpc server in order to satisfy envoy
		// health checks AND serve grpc requests, we can separate the http and
		// grpc port flags in k8s configs.
		sklog.Infof("regisering http healthz handler")
		topLevelRouter := mux.NewRouter()
		h := httputils.HealthzAndHTTPS(topLevelRouter)
		http.Handle("/", h)
		sklog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	}()

	s := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor))

	sklog.Infof("registering grpc health server")
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	sklog.Infof("registering grpc reflection server")
	reflection.Register(s)

	sklog.Infof("dialing RBE-CAS backends")
	rbeClients, err := backends.DialRBECAS(ctx)
	if err != nil {
		sklog.Fatalf("dialing RBE-CAS backends: %v", err)
	}
	sklog.Infof("successfully dialed %d RBE-CAS instances", len(rbeClients))

	sklog.Infof("dialing BigQuery")
	bqClient, err := backends.DialBigQuery(ctx)
	if err != nil {
		sklog.Fatalf("dialing bigquery: %v", err)
	}
	sklog.Infof("successfully dialed bigquery")

	sklog.Infof("registering cabe grpc server")
	cabeServer := analysisserver.New(rbeClients, bqClient)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		sklog.Fatalf("failed to listen: %v", err)
	}
	cpb.RegisterAnalysisServer(s, cabeServer)

	sklog.Infof("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		sklog.Fatalf("failed to serve grpc: %v", err)
	}
}
