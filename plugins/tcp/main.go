package main

import (
	"tcp/listener"

	"github.com/0xPrimo/TinyC2/sdk"
)

func NewListener(engine sdk.IEngine, name string, config string) (sdk.IListener, error) {
	return listener.NewTcpListener(engine, name, config)
}
