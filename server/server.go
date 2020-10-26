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
	"path/filepath"
	"strings"
)

func GRPCServer(ctx context.Context, cfg *config.Config, db *sql.DB) {
	// waitGroup add
	wg := config.GetWaitGroupInCtx(ctx)
	wg.Add(1)
	defer wg.Done()

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
	// waitGroup add
	wg := config.GetWaitGroupInCtx(ctx)
	wg.Add(1)
	defer wg.Done()

	// New Mux
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endpoint := strings.Join([]string{cfg.Host, cfg.GRPCPort}, ":")

	// registry handler endpoint
	err := v1.RegisterToDoServiceHandlerFromEndpoint(ctx, gwmux, endpoint, opts)
	if err != nil {
		log.Fatalf("registry http server error:%v", err)
	}

	// add swagger api
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	dir := filepath.Join(cfg.BaseDir, "api/proto/swagger")
	mux.Handle("/api/", http.StripPrefix("/api/", http.FileServer(http.Dir(dir))))

	log.Println("starting http server...")
	err = http.ListenAndServe(":"+cfg.HttpPort, mux)
	if err != nil {
		log.Fatalf("http server listen error:%v", err)
	}

	<-ctx.Done()
}
