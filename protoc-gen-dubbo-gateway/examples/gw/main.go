package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/extend"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-dubbo-gateway/examples"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()
	gw := gwruntime.NewServeMux()

	gatewayConfig := extend.NewDubboGatewayConfig(&extend.DubboGatewayOps{
		IsDirect: false,
		Protocol: "tri",
	})
	gatewayConfig.AddReferenceEndpoint("api", "tri://127.0.0.1:20000")
	gatewayConfig.AddRegistry("zookeeper", "127.0.0.1:2181", "127.0.0.1:2182")

	err := api.RegisterGreeterHandler(ctx, gw, gatewayConfig)
	if err != nil {
		panic(err)
	}

	err = api.RegisterGreeter2Handler(ctx, gw, gatewayConfig)
	if err != nil {
		panic(err)
	}

	err = gatewayConfig.Load()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gw)

	s := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	go func() {
		<-ctx.Done()
		glog.Infof("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			glog.Errorf("Failed to shutdown http server: %v", err)
		}
	}()

	glog.Infof("Starting listening at %s", "127.0.0.1:8080")
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		glog.Errorf("Failed to listen and serve: %v", err)
	}
}
