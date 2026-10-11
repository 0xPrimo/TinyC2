package engine

import (
	"context"
	"os"
	"os/exec"

	"github.com/0xPrimo/TinyC2/server/internal/database"
	"github.com/0xPrimo/TinyC2/server/internal/implant"
	"github.com/0xPrimo/TinyC2/server/internal/listener"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
	"github.com/0xPrimo/TinyC2/server/internal/plugin"
	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"
)

type Engine struct {
	db     *database.Database
	config *Config

	plugin.IPluginManager
	listener.IListenerManager
	implant.IImplantManager
}

func NewEngine(path string) (*Engine, error) {
	cfg, err := loadConfig(path)
	if err != nil {
		return nil, err
	}

	db, err := database.New("data.db")
	if err != nil {
		logger.Error("failed to initialize db: %v", err)
		return nil, err
	}

	err = db.Init(context.Background())
	if err != nil {
		return nil, err
	}

	pluginManager := plugin.NewManager(db)
	listenerManager := listener.NewManager(db, pluginManager)
	implantManager := implant.NewManager(db, listenerManager)

	return &Engine{
		config:           cfg,
		IPluginManager:   pluginManager,
		IListenerManager: listenerManager,
		IImplantManager:  implantManager,
	}, err
}

func (e *Engine) Init() error {
	// start cpl server
	cmd := exec.Command("cpl", "server")
	err := cmd.Start()
	if err != nil {
		logger.Error("failed to start cpl server: %v", err)
		return nil
	}
	logger.Success("cpl server started: http://127.0.0.1:60060/link")

	for _, pulg := range e.config.Plugins {
		meta, err := e.PluginRegister(e, pulg.Path)
		if err != nil {
			logger.Error("failed to register plugin: %v", err)
			return nil
		}

		logger.Success("plugin %s (%s) registered", pterm.Green(meta.Name), meta.Type)
	}

	e.PluginDBSync()
	e.ListenerDBSync()
	e.ImplantDBSync()
	return nil
}

func loadConfig(path string) (*Config, error) {
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
