package listener

import "github.com/0xPrimo/TinyC2/sdk"

type Listener struct {
	ID       uint32
	Name     string
	Protocol string
	sdk.IAdapterListener
}

type Meta struct {
	ID       uint32
	Name     string
	Protocol string
}
