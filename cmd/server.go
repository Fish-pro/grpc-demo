package cmd

import (
	"context"
	"database/sql"
	"fmt"
	v1 "github.com/Fish-pro/grpc-demo/api/service/v1"
	"github.com/Fish-pro/grpc-demo/config"
	"github.com/Fish-pro/grpc-demo/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func RunServer() error {
	cfg := config.New()
	ctx := context.Background()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server:%s", cfg.GRPCPort)
	}

	param := "parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.DbSchema, param)
	db, err := sql.Open("mysql", dsn)
	fmt.Println(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to mysql:%s", err.Error())
	}
	defer db.Close()

	v1API := v1.NewToDoServiceServer(db)

	go server.GRPCServer(ctx, v1API, cfg.GRPCPort)

	go server.HttpServer(ctx, cfg)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown server ...")

	return nil
}
