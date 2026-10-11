package listener

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"tcp/pkg/packer"

	"github.com/0xPrimo/TinyC2/sdk"

	"github.com/goccy/go-yaml"
)

type TCPListener struct {
	Name   string
	Engine sdk.IEngine
	Config TCPConfig

	listener net.Listener
	wg       sync.WaitGroup
	quit     chan struct{}
	stopOnce sync.Once
}

type Host struct {
	IP   string `yaml:"ip"`
	Port uint16 `yaml:"port"`
}

type TCPConfig struct {
	BindHost string `yaml:"bindhost"`
	BindPort string `yaml:"bindport"`
	Host     Host   `yaml:"host"`
}

func NewTCPListener(engine sdk.IEngine, name string, config string) (sdk.IListener, error) {
	data, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var cfg TCPConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &TCPListener{
		Name:   name,
		Engine: engine,
		Config: cfg,
	}, nil
}

func (t *TCPListener) Start() error {
	addr := t.Config.BindHost + ":" + t.Config.BindPort
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	t.listener = ln

	t.wg.Add(1)
	go t.acceptLoop()

	return nil
}

func (t *TCPListener) acceptLoop() {
	defer t.wg.Done()

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			// Handle other unexpected accept errors
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		t.wg.Add(1)
		go TCPHandler(t, conn)
	}
}

func (t *TCPListener) Stop() error {
	t.listener.Close()
	return nil
}

func (t *TCPListener) MakePic(id uint32) ([]byte, []byte, error) {
	configArray := []any{}

	configArray = append(configArray, uint32(id))
	configArray = append(configArray, t.Config.Host.IP)
	configArray = append(configArray, uint16(t.Config.Host.Port))

	piccfg, err := packer.Pack(configArray...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pack configurations: %v", err)
	}

	pic, err := os.ReadFile("plugins/tcp/bin/channel.x64.pic")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read pic: %w", err)
	}

	return pic, piccfg, nil
}

func IsValidAddress(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	_ = listener.Close()
	return nil
}
