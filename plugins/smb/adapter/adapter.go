package adapter

import (
	"fmt"
	"os"
	"smb/pkg/packer"

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

func (l *Listener) Extension(id uint32) ([]byte, []byte, error) {
	configArray := []any{
		uint32(id),
		l.config.PipeName,
	}

	piccfg, err := packer.Pack(configArray...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pack configurations: %v", err)
	}

	pic, err := os.ReadFile("../plugins/smb/build/smb.ext")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read extension: %w", err)
	}

	return pic, piccfg, nil
}

func (l *Listener) Config() map[string]any {
	return map[string]any{
		"pipename": l.config.PipeName,
	}
}
