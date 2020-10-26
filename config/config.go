package config

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	ctxKeyWaitGroup = "waitGroup"
)

type Config struct {
	Ctx      context.Context
	Cancel   context.CancelFunc
	Host     string
	GRPCPort string
	HttpPort string
	BaseDir  string
	OpenPem  bool
	Db       *Db
	Cert     *Certificate
}

type Db struct {
	Host     string
	User     string
	Password string
	DbSchema string
}

type Certificate struct {
	CaPath  string
	PemPath string
	KeyPath string
}

func GetWaitGroupInCtx(ctx context.Context) *sync.WaitGroup {
	if wg, ok := ctx.Value(ctxKeyWaitGroup).(*sync.WaitGroup); ok {
		return wg
	}

	return nil
}

func getBaseDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get baseDir path error: %v", err)
	}
	return currentDir
}

func toBool(val string) bool {
	return strings.ToLower(val) == "true"
}

func getEnvOrDefault(key string, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	} else {
		return def
	}
}

func New() *Config {
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), ctxKeyWaitGroup, new(sync.WaitGroup)))

	return &Config{
		Ctx:      ctx,
		Cancel:   cancel,
		BaseDir:  getBaseDir(),
		Host:     getEnvOrDefault("RUN_HOST", "localhost"),
		GRPCPort: getEnvOrDefault("GRPC_PORT", "8081"),
		HttpPort: getEnvOrDefault("GRPC_PORT", "8080"),
		OpenPem:  toBool(getEnvOrDefault("OPEN_PEM", "true")),
		Db: &Db{
			Host:     getEnvOrDefault("GRPC_HOST", "127.0.0.1:3306"),
			User:     getEnvOrDefault("GRPC_DB_USER", "root"),
			Password: getEnvOrDefault("GRPC_DB_PASSWORD", "dangerous"),
			DbSchema: getEnvOrDefault("GRPC_DB_SCHEMA", "grpc-demo"),
		},
		Cert: &Certificate{
			CaPath:  getEnvOrDefault("CA_PATH", "/Users/york/go/src/github.com/Fish-pro/grpc-demo/cert/ca.pem"),
			PemPath: getEnvOrDefault("PEM_PATH", "/Users/york/go/src/github.com/Fish-pro/grpc-demo/cert/server.pem"),
			KeyPath: getEnvOrDefault("KEY_PATH", "/Users/york/go/src/github.com/Fish-pro/grpc-demo/cert/server.key"),
		},
	}
}
