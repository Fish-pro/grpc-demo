package server

import (
	"context"
	v1 "github.com/Fish-pro/grpc-demo/api/proto/v1"
	"github.com/Fish-pro/grpc-demo/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
)

func GRPCServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("gRPC listen error:%v", err)
	}
	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, v1API)

	log.Println("starting gRPC server...")
	err = server.Serve(listen)
	if err != nil {
		log.Fatalf("start gRPC server error:%v", err)
	}

	<-ctx.Done()
}

func HttpServer(ctx context.Context, cfg *config.Config) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endpoint := strings.Join([]string{cfg.Host, cfg.GRPCPort}, ":")

	err := v1.RegisterToDoServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalf("registry http server error:%v", err)
	}

	log.Println("starting http server...")
	err = http.ListenAndServe(":"+cfg.HttpPort, mux)
	if err != nil {
		log.Fatalf("http server listen error:%v", err)
	}

	<-ctx.Done()
}
