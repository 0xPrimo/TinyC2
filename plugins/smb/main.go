package main

import (
	"smb/listener"

	"github.com/0xPrimo/TinyC2/sdk"
)

func NewListener(engine sdk.IEngine, name string, config string) (sdk.IListener, error) {
	return listener.NewSMBListener(engine, name, config)
}
