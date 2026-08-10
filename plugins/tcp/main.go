package main

import (
	"tcp/adapter"

	"github.com/0xPrimo/TinyC2/sdk"
)

var Plugin = PluginListener{}

type PluginListener struct {
	engine sdk.IEngine
}

func (p *PluginListener) Initialize(engine sdk.IEngine) {
	p.engine = engine
}

func (p *PluginListener) Meta() map[string]string {
	return map[string]string{
		"name": "tcp",
		"type": "listener",
	}
}

func (p *PluginListener) NewAdapter() sdk.IAdapterListener {
	return &adapter.Listener{
		IEngine: p.engine,
	}
}
