package handlers

import (
	"github.com/0xPrimo/TinyC2/server/internal/core"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"

	"github.com/pterm/pterm"
)

func HandleListener(engine *core.Engine, args []string) {
	if len(args) < 1 {
		logger.Info("listener [start|list|remove]")
		return
	}

	subcmd := args[0]
	switch subcmd {
	case "start":
		HandleListenerStart(engine, args[1:])
	case "stop":
		HandleListenerStop(engine, args[1:])
	case "generate":
		HandleListenerGenerate(engine, args[1:])
	case "list":
		HandleListenerList(engine, args[1:])
	}
}

func HandleListenerGenerate(engine *core.Engine, args []string) {
	if len(args) < 2 {
		logger.Info("listener generate [name] [dest]")
		return
	}

	name := args[0]
	dest := args[1]
	err := engine.ListenerGenerate(name, dest)
	if err != nil {
		logger.Error("listener generate %s: %v", name, err)
		return
	}

	logger.Success("implant generated successfully: %s", pterm.Cyan(dest))
}

func HandleListenerStart(engine *core.Engine, args []string) {
	if len(args) < 3 {
		logger.Info("listener start [plugin] [name] [config.yaml]")
		return
	}

	plugin := args[0]
	name := args[1]
	config := args[2]
	err := engine.ListenerStart(plugin, name, config)
	if err != nil {
		logger.Error("listener start %s: %v", plugin, err)
		return
	}

	logger.Success("listener %s started successfully", pterm.Cyan(name))
}

func HandleListenerList(engine *core.Engine, args []string) {
	err := engine.ListenerList()
	if err != nil {
		logger.Error("listener list: %v", err)
		return
	}
}

func HandleListenerStop(engine *core.Engine, args []string) {
	if len(args) < 1 {
		logger.Info("listener stop [name]")
		return
	}

	name := args[0]
	err := engine.ListenerStop(name)
	if err != nil {
		logger.Error("listener stop %s: %v", name, err)
		return
	}

	logger.Success("listener %s stopped successfully", pterm.Cyan(name))
}
