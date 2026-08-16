package adapter

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
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
	l.name = name
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

func (l *Listener) Extension(id uint32) ([]byte, error) {
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

	extcfg, err := packer.Pack(configArray...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack configurations: %w", err)
	}

	extdir, _ := filepath.Abs("../plugins/http/extension")
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
