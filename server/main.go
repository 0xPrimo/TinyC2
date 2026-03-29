package main

import (
	"os"

	"github.com/0xPrimo/TinyC2/server/internal/cli"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
)

func main() {
	if len(os.Args) < 2 {
		logger.Info("Usage: ./tinyc2 [config.yaml]")
		return
	}

	cli := cli.NewCli(os.Args[1])
	cli.Start()
}
