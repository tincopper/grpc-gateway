package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/config"
	"flag"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-dubbo-gateway/examples"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()
	gw := gwruntime.NewServeMux()

	refConf := config.ReferenceConfig{
		Protocol: "tri",
		URL:      "tri://127.0.0.1:20000",
	}

	err := api.RegisterGreeterHandler(ctx, gw, refConf)
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
