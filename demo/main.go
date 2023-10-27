package main

import (
	"fmt"

	"github.com/wxc/micro/config"
	"go-micro.dev/v4/logger"
)

func main() {
	if err := config.Load(); err != nil {
		logger.Fatal(err)
	}
	fmt.Printf("%s\n", config.Address())
}
