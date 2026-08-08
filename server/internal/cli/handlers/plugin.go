package handlers

import (
	"github.com/0xPrimo/TinyC2/server/internal/core"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"

	"github.com/pterm/pterm"
)

func HandlePlugin(engine *core.Engine, args []string) {
	if len(args) < 1 {
		logger.Info("plugin [register|list|unregister]")
		return
	}

	subcmd := args[0]
	switch subcmd {
	case "register":
		HandlePluginRegister(engine, args[1:])
	case "list":
		HandlePluginList(engine, args[1:])
	case "unregister":
		HandlePluginUnregister(engine, args[1:])
	}
}

func HandlePluginRegister(engine *core.Engine, args []string) {
	if len(args) < 2 {
		logger.Info("plugin register [path]")
		return
	}

	path := args[0]
	meta, err := engine.PluginRegister(engine, path)
	if err != nil {
		logger.Error("plugin load: %v", err)
		return
	}

	logger.Success("plugin %s ( %s ) registered successfully", pterm.Green(meta.Name), pterm.Cyan(meta.Type))
}

func HandlePluginList(engine *core.Engine, args []string) {

	table := pterm.TableData{
		{"Name", "Type", "Path"},
	}

	pls := engine.PluginList()
	for _, pl := range pls {
		table = append(table, []string{pterm.Cyan(pl.Name), pl.Type, pl.Path})

	}
	pterm.Println()
	pterm.DefaultTable.
		WithHasHeader().
		WithBoxed().
		WithHeaderStyle(pterm.NewStyle(pterm.FgLightMagenta, pterm.Bold)).
		WithData(table).
		Render()
	pterm.Println()
}

func HandlePluginUnregister(engine *core.Engine, args []string) {
	if len(args) < 1 {
		logger.Info("plugin remove [name]")
		return
	}

	name := args[0]
	err := engine.PluginUnregister(name)
	if err != nil {
		logger.Error("plugin remove: %v", err)
		return
	}

	logger.Success("plugin %s removed successfully", pterm.Green(name))
}
