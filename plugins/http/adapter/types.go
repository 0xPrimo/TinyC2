package adapter

const (
	CPL_SERVER = "http://127.0.0.1:60060/link"
	SPEC_FILE  = "channel.spec"
	CAPAB_FILE = "channel.x64.o"
)

type CplResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	YaraRule  string `json:"yara"`
	OutputB64 string `json:"output_b64"`
	Context   string `json:"context"`
}

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
