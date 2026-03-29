package handlers

import (
	"github.com/0xPrimo/TinyC2/server/internal/core"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"

	"github.com/pterm/pterm"
)

func HandlePlugin(engine *core.Engine, args []string) {
	if len(args) < 1 {
		logger.Info("plugin [register|list|remove]")
		return
	}

	subcmd := args[0]
	switch subcmd {
	case "register":
		HandlePluginRegister(engine, args[1:])
	case "list":
		HandlePluginList(engine, args[1:])
	case "remove":
		HandlePluginList(engine, args[1:])
	}
}

func HandlePluginRegister(engine *core.Engine, args []string) {
	if len(args) < 2 {
		logger.Info("plugin register [name] [config.yaml]")
		return
	}

	name := args[0]
	config := args[1]
	err := engine.PluginRegister(name, config)
	if err != nil {
		logger.Error("plugin load: %v", err)
		return
	}

	logger.Success("plugin %s registred successfully", pterm.Green(name))
}

func HandlePluginList(engine *core.Engine, args []string) {
	err := engine.PluginList()
	if err != nil {
		logger.Error("plugin list: %v", err)
		return
	}
}

func HandlePluginRemove(engine *core.Engine, args []string) {
	if len(args) < 1 {
		logger.Info("plugin remove [name]")
		return
	}

	name := args[0]
	err := engine.PluginRemove(name)
	if err != nil {
		logger.Error("plugin remove: %v", err)
		return
	}

	logger.Success("plugin %s removed successfully", pterm.Green(name))
}
