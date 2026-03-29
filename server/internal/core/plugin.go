package core

import (
	"fmt"
	"plugin"

	"github.com/0xPrimo/TinyC2/sdk"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"

	"github.com/pterm/pterm"
)

type Plugin struct {
	Name        string
	Path        string
	NewListener func(engine sdk.IEngine, name string, config string) (sdk.IListener, error)
}

// PluginRegister register plugin
func (e *Engine) PluginRegister(name string, path string) error {
	plug, err := plugin.Open(path)
	if err != nil {
		return err
	}

	sym, err := plug.Lookup("NewListener")
	if err != nil {
		return err
	}

	NewListener, ok := sym.(func(sdk.IEngine, string, string) (sdk.IListener, error))
	if !ok {
		return fmt.Errorf("symbol 'NewListener' has the wrong signature")
	}

	e.Plugins[name] = Plugin{
		Name:        name,
		Path:        path,
		NewListener: NewListener,
	}

	logger.Success("plugin %s registred successfully", pterm.Green(name))

	return nil
}

// PluginRemove remove plugin
func (e *Engine) PluginRemove(name string) error {
	delete(e.Plugins, name)
	return nil
}

// PluginList list plugins
func (e *Engine) PluginList() error {
	table := pterm.TableData{
		{"Name", "Path"},
	}

	for _, plugin := range e.Plugins {
		table = append(table, []string{pterm.Cyan(plugin.Name), plugin.Path})
	}

	pterm.Println()
	pterm.DefaultTable.
		WithHasHeader().
		WithBoxed().
		WithHeaderStyle(pterm.NewStyle(pterm.FgLightMagenta, pterm.Bold)).
		WithData(table).
		Render()
	pterm.Println()

	return nil
}

// PluginNewListener create new listener using the plugin
func (e *Engine) PluginNewListener(plugin string, name string, config string) (sdk.IListener, error) {
	plug, exists := e.Plugins[plugin]
	if !exists {
		return nil, fmt.Errorf("plugin %s does not exist", name)
	}

	return plug.NewListener(e, name, config)
}
