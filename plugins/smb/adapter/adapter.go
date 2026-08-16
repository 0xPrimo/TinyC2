package adapter

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"smb/pkg/packer"
	"time"

	"github.com/0xPrimo/TinyC2/sdk"
	"github.com/goccy/go-yaml"
)

type Listener struct {
	name   string
	config Config
	sdk.IEngine
}

func (l *Listener) Start(name string, config string) error {
	data, err := os.ReadFile(config)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &l.config); err != nil {
		return err
	}

	return nil
}

func (l *Listener) Stop() error {
	return nil
}

func (l *Listener) Extension(id uint32) ([]byte, error) {
	configArray := []any{
		uint32(id),
		l.config.PipeName,
	}

	extcfg, err := packer.Pack(configArray...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack configurations: %w", err)
	}

	extdir, _ := filepath.Abs("../plugins/smb/extension")
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
		"pipename": l.config.PipeName,
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
