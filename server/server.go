package server

import (
	"context"
	"database/sql"
	v1 "github.com/Fish-pro/grpc-demo/api/proto/v1"
	serviceV1 "github.com/Fish-pro/grpc-demo/api/service/v1"
	"github.com/Fish-pro/grpc-demo/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
)

func GRPCServer(ctx context.Context, db *sql.DB, cfg *config.Config) {
	// New tcp listen
	listen, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("gRPC listen error:%v", err)
	}

	// New grpc server
	server := grpc.NewServer()
	// registry todoServiceServer in version v1
	v1ToDoAPI := serviceV1.NewToDoServiceServer(db)
	v1.RegisterToDoServiceServer(server, v1ToDoAPI)

	log.Println("starting gRPC server...")
	err = server.Serve(listen)
	if err != nil {
		log.Fatalf("start gRPC server error:%v", err)
	}

	<-ctx.Done()
}

func HttpServer(ctx context.Context, cfg *config.Config) {
	// New Mux
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endpoint := strings.Join([]string{cfg.Host, cfg.GRPCPort}, ":")

	// registry handler endpoint
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
