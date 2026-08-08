package plug

import "github.com/0xPrimo/TinyC2/sdk"

type PluginListener struct {
	Name string
	Path string

	sdk.IPluginListener
}

type Meta struct {
	Name   string
	Type   string
	Config string
	Path   string
}
