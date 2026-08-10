package adapter

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"tcp/pkg/packer"

	"github.com/0xPrimo/TinyC2/sdk"
	"github.com/goccy/go-yaml"
)

type Listener struct {
	name   string
	config Config

	listener net.Listener
	wg       sync.WaitGroup
	stopOnce sync.Once

	conns   map[net.Conn]struct{}
	connsMu sync.Mutex

	sdk.IEngine
}

func (l *Listener) Start(name string, configPath string) error {
	l.name = name

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %q: %w", configPath, err)
	}

	if err := yaml.Unmarshal(data, &l.config); err != nil {
		return fmt.Errorf("failed to parse YAML configuration: %w", err)
	}

	addr := net.JoinHostPort(l.config.BindHost, l.config.BindPort)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	l.listener = ln
	l.conns = make(map[net.Conn]struct{})

	l.LogInfo("tcp", "Started listener %s on %s", name, addr)

	l.wg.Add(1)
	go func() {
		defer l.wg.Done()

		for {
			conn, err := l.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				l.LogError("tcp", "Error accepting connection: %v", err)
				continue
			}

			l.connsMu.Lock()
			if l.conns == nil {
				conn.Close()
				l.connsMu.Unlock()
				return
			}
			l.conns[conn] = struct{}{}
			l.connsMu.Unlock()

			// pass the connection to handler
			l.wg.Add(1)
			go l.handler(conn)
		}
	}()

	return nil
}

func (l *Listener) Stop() error {
	var err error
	l.stopOnce.Do(func() {
		if l.listener != nil {
			err = l.listener.Close()
		}

		// Close all active connections to unblock read/write loops
		l.connsMu.Lock()
		for conn := range l.conns {
			conn.Close()
		}
		l.conns = nil
		l.connsMu.Unlock()

		// Wait for accept loop, handlers, and writer goroutines to exit
		l.wg.Wait()
		l.LogInfo("tcp", "Listener %s stopped", l.name)
	})
	return err
}

func (l *Listener) Extension(id uint32) ([]byte, []byte, error) {
	configArray := []any{
		uint32(id),
		l.config.Host.IP,
		uint16(l.config.Host.Port),
	}

	piccfg, err := packer.Pack(configArray...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pack configurations: %w", err)
	}

	// Read extension binary
	extensionPath := "../plugins/tcp/build/tcp.ext"
	pic, err := os.ReadFile(extensionPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read extension from %q: %w", extensionPath, err)
	}

	return pic, piccfg, nil
}

func (l *Listener) Config() map[string]any {
	return map[string]any{
		"bindhost": l.config.BindHost,
		"bindport": l.config.BindPort,
		"host": map[string]any{
			"ip":   l.config.Host.IP,
			"port": l.config.Host.Port,
		},
	}
}
