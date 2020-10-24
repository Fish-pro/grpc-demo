package config

import (
	"log"
	"os"
)

type Config struct {
	Host     string
	GRPCPort string
	HttpPort string
	BaseDir  string
	Db       *Db
}

type Db struct {
	Host     string
	User     string
	Password string
	DbSchema string
}

func getBaseDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current path error: %v", err)
	}
	return currentDir
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
		BaseDir:  getBaseDir(),
		Host:     getEnvOrDefault("RUN_HOST", "localhost"),
		GRPCPort: getEnvOrDefault("GRPC_PORT", "8081"),
		HttpPort: getEnvOrDefault("GRPC_PORT", "8080"),
		Db: &Db{
			Host:     getEnvOrDefault("GRPC_HOST", "127.0.0.1:3306"),
			User:     getEnvOrDefault("GRPC_DB_USER", "root"),
			Password: getEnvOrDefault("GRPC_DB_PASSWORD", "dangerous"),
			DbSchema: getEnvOrDefault("GRPC_DB_SCHEMA", "grpc-demo"),
		},
	}
}
