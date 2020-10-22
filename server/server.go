package server

import (
	"context"
	v1 "github.com/Fish-pro/grpc-demo/api/proto/v1"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	go func() {
		for range c {
			log.Println("shutting dow gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
