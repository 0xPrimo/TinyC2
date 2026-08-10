package adapter

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

type Listener struct {
	name   string
	server *http.Server
	config Config

	sdk.IEngine
}

func (l *Listener) Start(name string, configPath string) error {

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &l.config); err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", l.config.BindHost, l.config.BindPort)
	if err := isValidAddress(addr); err != nil {
		return err
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.POST("/*any", l.handler())

	l.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := l.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.LogError(l.name, "listen and serve error: %v", err)
		}
	}()

	return nil
}

func (l *Listener) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := l.server.Shutdown(ctx); err != nil {
		l.LogWarning(l.name, "listener forced to shutdown: %v", err)
	}

	return nil
}

func (l *Listener) Extension(id uint32) ([]byte, []byte, error) {
	configArray := []any{
		uint32(id),
		l.config.UserAgent,
		l.config.Method,
	}

	switch l.config.Rotation {
	case "round-robin":
		configArray = append(configArray, uint16(1))
	default:
		configArray = append(configArray, uint16(0))
	}

	configArray = append(configArray, uint16(len(l.config.Uris)))
	configArray = append(configArray, uint16(len(l.config.Headers)))
	configArray = append(configArray, uint16(len(l.config.Hosts)))

	for _, host := range l.config.Hosts {
		configArray = append(configArray, host.IP)
		configArray = append(configArray, uint16(host.Port))
	}

	for key, value := range l.config.Headers {
		configArray = append(configArray, key+": "+value+"\r\n")
	}

	for _, uri := range l.config.Uris {
		configArray = append(configArray, uri)
	}

	picargs, err := packer.Pack(configArray...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pack configurations: %v", err)
	}

	data, err := os.ReadFile("../plugins/http/build/http.ext")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read extension: %w", err)
	}

	return data, picargs, nil
}

func (l *Listener) Config() map[string]any {
	return map[string]any{
		"bindhost":   l.config.BindHost,
		"bindport":   l.config.BindPort,
		"hosts":      l.config.Hosts,
		"user-agent": l.config.UserAgent,
		"method":     l.config.Method,
	}
}

func isValidAddress(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	_ = listener.Close()
	return nil
}
