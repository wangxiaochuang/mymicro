package main

import (
	"fmt"

	"github.com/wxc/demo/config"
	"github.com/wxc/micro/logger"
)

func main() {
	if err := config.Load(); err != nil {
		logger.Fatal(err)
	}
	fmt.Printf("%s %s", config.Address(), config.Tracing().Jaeger.URL)
}
