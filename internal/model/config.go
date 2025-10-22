package model

type Config struct {
	Token    string `yaml:"token"`
	Endpoint string `yaml:"endpoint"`
	CaCert   string `yaml:"ca_cert,omitempty"`
}
