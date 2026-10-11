package adapter

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"tcp/pkg/packer"
	"time"

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

func (l *Listener) Extension(id uint32) ([]byte, error) {
	configArray := []any{
		uint32(id),
		l.config.Host.IP,
		uint16(l.config.Host.Port),
	}

	extcfg, err := packer.Pack(configArray...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack configurations: %w", err)
	}

	extdir, _ := filepath.Abs("../plugins/tcp/extension")
	resp, err := SendHttp[CplResponse](CPL_SERVER, map[string]any{
		"action": "link",
		"params": map[string]any{
			"spec": filepath.Join(extdir, SPEC_FILE),
			"file": filepath.Join(extdir, "bin", CAPAB_FILE),
		},
		"env": map[string]any{
			"$CONFIG": hex.EncodeToString(extcfg),
		},
	})
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to build extension: \n\tmessage: %s\n\tcontext: %s\n\t", resp.Message, resp.Context)
	}

	pic, _ := base64.StdEncoding.DecodeString(resp.OutputB64)
	return pic, nil
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

func SendHttp[T any](url string, body map[string]any) (*T, error) {
	var result T

	reqbody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqbody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
