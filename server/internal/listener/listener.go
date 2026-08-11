package listener

import (
	"fmt"
	"hash/crc32"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/store"
	"github.com/0xPrimo/TinyC2/server/internal/plug"
)

type IListenerManager interface {
	ListenerStart(plugin string, name string, config string) error
	ListenerStop(name string) error
	ListenerExtension(name string) ([]byte, []byte, error)
	ListenerConfig(name string) (map[string]any, error)
	ListenerList() []Meta
}

type Manager struct {
	listeners *store.Store[string, *Listener]
	plug.IPluginManager
}

func NewManager(pluginManager plug.IPluginManager) *Manager {
	return &Manager{
		IPluginManager: pluginManager,
		listeners:      store.NewStore[string, *Listener](),
	}
}

func (m *Manager) ListenerStart(plugin string, name string, config string) error {
	if m.listeners.Has(name) {
		return fmt.Errorf("listener %s already exists", name)
	}

	pl, ok := m.PluginListener(plugin)
	if !ok {
		return fmt.Errorf("unknown plugin %s", plugin)
	}

	listener := &Listener{
		ID:               crc32.ChecksumIEEE([]byte(name)),
		Name:             name,
		Protocol:         plugin,
		IAdapterListener: pl.NewAdapter(),
	}

	err := listener.Start(name, config)
	if err != nil {
		return err
	}

	m.listeners.Set(name, listener)

	return nil
}

func (m *Manager) ListenerStop(name string) error {
	listener, ok := m.listeners.Get(name)
	if !ok {
		return fmt.Errorf("listener %s doesn't exists", name)
	}

	if err := listener.Stop(); err != nil {
		return err
	}

	m.listeners.Delete(name)

	return nil

}
func (m *Manager) ListenerExtension(name string) ([]byte, []byte, error) {
	listener, ok := m.listeners.Get(name)
	if !ok {
		return nil, nil, fmt.Errorf("listener %s doesn't exists", name)
	}

	pic, piccfg, err := listener.Extension(0)
	if err != nil {
		return nil, nil, err
	}

	return pic, piccfg, nil
}

func (m *Manager) ListenerConfig(name string) (map[string]any, error) {
	listener, ok := m.listeners.Get(name)
	if !ok {
		return nil, fmt.Errorf("listener %s doesn't exists", name)
	}

	return listener.Config(), nil
}

func (m *Manager) ListenerList() []Meta {
	var lns []Meta

	m.listeners.ForEach(func(name string, listener *Listener) {
		lns = append(lns, Meta{
			ID:       listener.ID,
			Name:     listener.Name,
			Protocol: listener.Protocol,
		})
	})

	return lns
}
