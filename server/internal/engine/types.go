package engine

import (
	"github.com/0xPrimo/TinyC2/server/internal/implant"
	"github.com/0xPrimo/TinyC2/server/internal/listener"
	"github.com/0xPrimo/TinyC2/server/internal/plugin"
)

type IEngine interface {
	implant.IImplantManager
	listener.IListenerManager
	plugin.IPluginManager
}

type Config struct {
	Plugins       []PluginConfig      `yaml:"plugins"`
	CrystalPalace CrystalPalaceConfig `yaml:"crystal-palace"`
}

type CrystalPalaceConfig struct {
	Lib      string `yaml:"lib"`
	Pavilion string `yaml:"pavilion"`
}

type PluginConfig struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}
