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

type Config struct {
	PipeName string `yaml:"pipename"`
}
