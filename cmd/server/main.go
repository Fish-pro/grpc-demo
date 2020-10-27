package main

import (
	"fmt"
	"github.com/Fish-pro/grpc-demo/cmd"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
