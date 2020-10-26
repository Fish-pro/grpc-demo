package cmd

import (
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
	ctx := cfg.Ctx
	defer func() {
		cfg.Cancel()
		config.GetWaitGroupInCtx(ctx).Wait() // wait for goroutine cancel
	}()

	// add database connect information
	db, err := config.InitDb(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	// grpc server
	go server.GRPCServer(ctx, cfg, db)

	// http server
	go server.HttpServer(ctx, cfg)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown server ...")

	return nil
}
