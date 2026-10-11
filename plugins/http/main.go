package main

import (
	"http/listener"

	"github.com/0xPrimo/TinyC2/sdk"
)

func NewListener(engine sdk.IEngine, name string, config string) (sdk.IListener, error) {
	return listener.NewHTTPListener(engine, name, config)
}
