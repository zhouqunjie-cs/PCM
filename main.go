package main

import (
	"context"
	"flag"
	"github.com/JCCE-nudt/PCM/common/server"
	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/demo"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbecs"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbpod"
	"net"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", ":9090", ":8081")
)

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := demo.RegisterDemoServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return errors.Wrap(err, "RegisterDemoServiceHandlerFromEndpoint error")
	} else if err = pbecs.RegisterEcsServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return errors.Wrap(err, "RegisterEcsServiceHandlerFromEndpoint error")
	} else if err = pbpod.RegisterPodServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return errors.Wrap(err, "RegisterPodServiceHandlerFromEndpoint error")
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", mux)
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "conf", "configs/config.yaml", "config.yaml")
	flag.Parse()
	defer glog.Flush()

	if err := tenanter.LoadCloudConfigsFromFile(configFile); err != nil {
		if !errors.Is(err, tenanter.ErrLoadTenanterFileEmpty) {
			glog.Fatalf("LoadCloudConfigsFromFile error %+v", err)
		}
		glog.Warningf("LoadCloudConfigsFromFile empty file path %s", configFile)
	}

	glog.Infof("load tenant from file finished")

	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			glog.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		demo.RegisterDemoServiceServer(s, &server.Server{})
		pbecs.RegisterEcsServiceServer(s, &server.Server{})
		pbpod.RegisterPodServiceServer(s, &server.Server{})

		if err = s.Serve(lis); err != nil {
			glog.Fatalf("failed to serve: %v", err)
		}
	}()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
