package core

import (
	"os"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"

	"gopkg.in/yaml.v3"
)

type Engine struct {
	// TODO: use refrence instead of value for Plugins, Listeners and Impalnts map
	//
	Plugins   map[string]Plugin
	Implants  map[uint32]Implant
	Listeners map[string]Listener
	Config    EngineConfig
}

type EngineConfig struct {
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

func NewEngine(path string) *Engine {
	var (
		config EngineConfig
	)

	data, err := os.ReadFile(path)
	if err != nil {
		logger.Error("error reading file: %v", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logger.Error("error unmarshaling YAML: %v", err)
		os.Exit(1)
	}

	logger.Info("config: %v", config)

	engine := &Engine{
		Plugins:   make(map[string]Plugin),
		Implants:  make(map[uint32]Implant),
		Listeners: make(map[string]Listener),
		Config:    config,
	}

	for _, plugin := range config.Plugins {
		err := engine.PluginRegister(plugin.Name, plugin.Path)
		if err != nil {
			logger.Error("failed to register plugin: %v", err)
		}
	}

	return engine
}
