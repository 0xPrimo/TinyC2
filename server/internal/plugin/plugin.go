// Package plugin plugin management
package plugin

import (
	"fmt"
	"plugin"

	"github.com/0xPrimo/TinyC2/sdk"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/store"
)

type IPluginManager interface {
	PluginRegister(engine sdk.IEngine, path string) (Meta, error)
	PluginUnregister(name string) error
	PluginList() []Meta
	PluginListener(name string) (*PluginListener, bool)
}

type Manager struct {
	PluginsMeta map[string]Meta
	Listeners   *store.Store[string, *PluginListener]
}

func NewManager() *Manager {
	return &Manager{
		Listeners: store.NewStore[string, *PluginListener](),
	}
}

func (m *Manager) PluginRegister(engine sdk.IEngine, path string) (Meta, error) {
	mod, err := plugin.Open(path)
	if err != nil {
		return Meta{}, err
	}

	sym, err := mod.Lookup("Plugin")
	if err != nil {
		return Meta{}, err
	}

	pl, ok := sym.(sdk.IPlugin)
	if !ok {
		return Meta{}, fmt.Errorf("export Plugin does not implement IPlugin")
	}

	meta := pl.Meta()
	switch meta["type"] {
	case "listener":

		pl.Initialize(engine)
		m.Listeners.Set(meta["name"], &PluginListener{
			IPluginListener: pl.(sdk.IPluginListener),
			Name:            meta["name"],
			Path:            path,
		})

		return Meta{Name: meta["name"], Type: meta["type"]}, nil
	default:
		return Meta{}, fmt.Errorf("unknown plugin type: %s", meta["type"])
	}
}

func (m *Manager) PluginUnregister(name string) error {
	meta, ok := m.PluginsMeta[name]
	if !ok {
		return fmt.Errorf("unknown plugin: %s", name)
	}

	switch meta.Type {
	case "listener":
		m.Listeners.Delete(name)
	}

	delete(m.PluginsMeta, name)
	return nil
}

func (m *Manager) PluginListener(name string) (*PluginListener, bool) {
	pl, ok := m.Listeners.Get(name)
	if !ok {
		return nil, false
	}

	return pl, true
}

func (m *Manager) PluginList() []Meta {
	var pls []Meta

	m.Listeners.ForEach(func(name string, plugin *PluginListener) {
		pls = append(pls, Meta{
			Name: name,
			Path: plugin.Path,
			Type: "listener",
		})
	})

	return pls
}
