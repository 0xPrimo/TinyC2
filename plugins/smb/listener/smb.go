package listener

import (
	"fmt"
	"os"

	"smb/pkg/packer"

	"github.com/0xPrimo/TinyC2/sdk"

	"github.com/goccy/go-yaml"
)

type SMBListener struct {
	Name   string
	Engine sdk.IEngine
	Config SMBConfig
}

type Host struct {
	PipeName string `yaml:"pipename"`
}

type SMBConfig struct {
	PipeName string `yaml:"pipename"`
}

func NewSMBListener(engine sdk.IEngine, name string, config string) (sdk.IListener, error) {
	data, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var cfg SMBConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &SMBListener{
		Name:   name,
		Engine: engine,
		Config: cfg,
	}, nil
}

func (t *SMBListener) Start() error {
	return nil
}

func (t *SMBListener) Stop() error {
	return nil
}

func (t *SMBListener) MakePic(id uint32) ([]byte, []byte, error) {
	configArray := []any{}

	configArray = append(configArray, uint32(id))
	configArray = append(configArray, t.Config.PipeName)

	piccfg, err := packer.Pack(configArray...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pack configurations: %v", err)
	}

	pic, err := os.ReadFile("plugins/smb/bin/channel.x64.pic")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read pic: %w", err)
	}

	return pic, piccfg, nil
}
