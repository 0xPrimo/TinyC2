package adapter

type Host struct {
	IP   string `yaml:"ip"`
	Port uint16 `yaml:"port"`
}

type Config struct {
	BindHost string `yaml:"bindhost"`
	BindPort string `yaml:"bindport"`
	Host     Host   `yaml:"host"`
}
