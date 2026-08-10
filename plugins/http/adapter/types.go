package adapter

type Host struct {
	IP   string `yaml:"ip"`
	Port uint16 `yaml:"port"`
}

type Config struct {
	BindHost  string            `yaml:"bindhost"`
	BindPort  uint16            `yaml:"bindport"`
	Hosts     []Host            `yaml:"hosts"`
	Rotation  string            `yaml:"rotation-strategy"`
	UserAgent string            `yaml:"user-agent"`
	Method    string            `yaml:"method"`
	Uris      []string          `yaml:"uris"`
	Headers   map[string]string `yaml:"headers"`
}
