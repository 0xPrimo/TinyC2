package main

import (
	"protocol/listener"

	"github.com/0xPrimo/TinyC2/sdk"
)

func NewListener(engine sdk.IEngine, name string, config string) (sdk.IListener, error) {
	return listener.NewProtocolListener(engine, name, config)
}
