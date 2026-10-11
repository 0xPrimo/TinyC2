package listener

import (
	"errors"
	"fmt"
	"hash/crc32"

	"github.com/0xPrimo/TinyC2/server/internal/database"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/store"
	"github.com/0xPrimo/TinyC2/server/internal/plugin"
)

type IListenerManager interface {
	ListenerStart(plugin string, name string, config string) error
	ListenerStop(name string) error
	ListenerExtension(name string) ([]byte, error)
	ListenerConfig(name string) (map[string]any, error)
	ListenerList() []Meta
	ListenerGet(name string) (Meta, bool)
	ListenerGetByID(id uint32) (Meta, bool)
	ListenerDBSync() error
}

type Manager struct {
	db        *database.Database
	listeners *store.Store[string, *Listener]

	plugin.IPluginManager
}

func NewManager(db *database.Database, pluginManager plugin.IPluginManager) *Manager {
	return &Manager{
		db:             db,
		listeners:      store.NewStore[string, *Listener](),
		IPluginManager: pluginManager,
	}
}

func (m *Manager) ListenerDBSync() error {
	var errs []error

	listeners, err := m.db.ListenerGetAll()
	if err != nil {
		return err
	}

	for _, meta := range listeners {
		listener, err := m.createListener(meta.Protocol, meta.Name, meta.Config)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to start listener: %v", err))
			continue
		}

		err = listener.Start(meta.Name, meta.Config)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to start listener: %v", err))
		}
	}

	return errors.Join(errs...)
}

func (m *Manager) createListener(plugin, name, config string) (*Listener, error) {
	pl, ok := m.PluginListener(plugin)
	if !ok {
		return nil, fmt.Errorf("unknown plugin %s", plugin)
	}

	listener := &Listener{
		ID:               crc32.ChecksumIEEE([]byte(name)),
		Name:             name,
		Protocol:         plugin,
		IAdapterListener: pl.NewAdapter(),
	}

	m.listeners.Set(name, listener)

	return listener, nil
}

func (m *Manager) ListenerStart(plugin string, name string, config string) error {
	if m.listeners.Has(name) {
		return fmt.Errorf("listener %s already exists", name)
	}

	listener, err := m.createListener(plugin, name, config)
	if err != nil {
		return err
	}

	err = listener.Start(name, config)
	if err != nil {
		return err
	}

	err = m.db.ListenerCreate(database.Listener{listener.Name, listener.Protocol, config})
	if err != nil {
		return err
	}

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
func (m *Manager) ListenerExtension(name string) ([]byte, error) {
	listener, ok := m.listeners.Get(name)
	if !ok {
		return nil, fmt.Errorf("listener %s doesn't exists", name)
	}

	logger.Info("generating extension with id: %X", listener.ID)
	pic, err := listener.Extension(listener.ID)
	if err != nil {
		return nil, err
	}

	return pic, nil
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

func (m *Manager) ListenerGet(name string) (Meta, bool) {
	listener, ok := m.listeners.Get(name)
	if !ok {
		return Meta{}, false
	}

	return Meta{
		ID:       listener.ID,
		Name:     listener.Name,
		Protocol: listener.Protocol,
	}, true
}

func (m *Manager) ListenerGetByID(id uint32) (Meta, bool) {
	var ln Meta

	m.listeners.ForEach(func(name string, listener *Listener) {
		if listener.ID == id {
			ln.Name = listener.Name
			ln.Protocol = listener.Protocol
			ln.ID = listener.ID
		}
	})

	if ln.ID == 0 {
		return Meta{}, false
	}
	return ln, true
}
