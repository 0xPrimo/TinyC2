package listener

import (
	"os"

	"github.com/0xPrimo/TinyC2/sdk"

	"github.com/goccy/go-yaml"
)

type ProtocolListener struct {
	// Name
	Name string

	// Engine interface. logging and implant APIs
	Engine sdk.IEngine

	// Listener config
	Config ProtocolConfig
}

type ProtocolConfig struct {
	ConnHost string `yaml:"connhost"`
	ConnPort uint16 `yaml:"connport"`
	BindHost string `yaml:"bindhost"`
	BindPort uint16 `yaml:"bindport"`
}

// NewProtocolListener plugin manager will call this function when the user
// wants to create a new listener object. the function will receive engine interface,
// name and path to config as argument
func NewProtocolListener(engine sdk.IEngine, name string, config string) (sdk.IListener, error) {
	data, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var cfg ProtocolConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// This object should match IListener interface
	return &ProtocolListener{
		Name:   name,
		Engine: engine,
		Config: cfg,
	}, nil
}

// This will be called on command `listener start [plugin name] [name] [config.yaml]`
func (h *ProtocolListener) Start() error {
	return nil
}

// This will be called on command `listener stop [name]`
func (h *ProtocolListener) Stop() error {
	return nil
}

// This will be called on command `listener generate [name] [dest]` or
// `channel register [name]`
func (h *ProtocolListener) MakePic(id uint32) ([]byte, error) {
	// return Crystal Palace PIC
	return nil, nil
}
