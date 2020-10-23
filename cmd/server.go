package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Fish-pro/grpc-demo/config"
	"github.com/Fish-pro/grpc-demo/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func RunServer() error {
	// New Configuration information
	cfg := config.New()
	// New context information
	ctx := context.Background()

	// add database connect information
	param := "parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.DbSchema, param)
	db, err := sql.Open("mysql", dsn)
	log.Println("database info>>>", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to mysql:%s", err.Error())
	}
	defer db.Close()

	// grpc server
	go server.GRPCServer(ctx, db, cfg)

	// http server
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
