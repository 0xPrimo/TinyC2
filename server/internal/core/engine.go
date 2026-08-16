package core

import (
	"os"
	"os/exec"

	"github.com/0xPrimo/TinyC2/server/internal/listener"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
	"github.com/0xPrimo/TinyC2/server/internal/plug"
	"github.com/pterm/pterm"

	"gopkg.in/yaml.v3"
)

type Engine struct {
	Implants map[uint32]Implant
	Config   EngineConfig

	//Listeners *store.Store[string, *listener.Listener]
	plug.IPluginManager
	listener.IListenerManager
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
	var config EngineConfig

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

	pluginManager := plug.NewManager()
	listenerManager := listener.NewManager(pluginManager)
	engine := &Engine{
		Implants:         make(map[uint32]Implant),
		Config:           config,
		IPluginManager:   pluginManager,
		IListenerManager: listenerManager,
	}

	// start cpl server
	cmd := exec.Command("cpl", "server")
	err = cmd.Start()
	if err != nil {
		logger.Error("failed to start cpl server: %v", err)
		return nil
	}

	logger.Success("cpl server started: http://127.0.0.1:60060/link")

	for _, plugin := range config.Plugins {
		meta, err := engine.PluginRegister(engine, plugin.Path)
		if err != nil {
			logger.Error("failed to register plugin: %v", err)
			return nil
		}

		logger.Success("plugin %s (%s) registered", pterm.Green(meta.Name), meta.Type)
	}

	SetupImplantTaskResultHandlers()

	return engine
}
