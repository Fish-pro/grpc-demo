package cmd

import (
	"context"
	"database/sql"
	"fmt"
	v1 "github.com/Fish-pro/grpc-demo/api/service/v1"
	"github.com/Fish-pro/grpc-demo/server"
	"os"
)

type Config struct {
	GRPCPort string
	Db       *Db
}

type Db struct {
	Host     string
	User     string
	Password string
	DbSchema string
}

func getEnvOrDefault(key string, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	} else {
		return def
	}
}

func New() *Config {
	return &Config{
		GRPCPort: getEnvOrDefault("GRPC_PORT", "8080"),
		Db: &Db{
			Host:     getEnvOrDefault("GRPC_HOST", "127.0.0.1:3306"),
			User:     getEnvOrDefault("GRPC_DB_USER", "root"),
			Password: getEnvOrDefault("GRPC_DB_PASSWORD", "dangerous"),
			DbSchema: getEnvOrDefault("GRPC_DB_SCHEMA", "grpc-demo"),
		},
	}
}

func RunServer() error {
	ctx := context.Background()
	cfg := New()
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

	return server.RunServer(ctx, v1API, cfg.GRPCPort)
}
