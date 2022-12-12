package config

type Config struct {
	Input       string            `yaml:"input" validate:"required"`
	Output      string            `yaml:"output" validate:"required"`
	Marketplace MarketPlaceConfig `yaml:"marketplace" validate:"required"`
}

type MarketPlaceConfig struct {
	Appkey      string      `yaml:"appkey" validate:"omitempty"`
	Mktplaceurl string      `yaml:"mktplaceurl" validation:"omitempty"`
	Credentials Credentials `yaml:"credentials" validate:"omitempty"`
}

type Credentials struct {
	Bearer       string       `yaml:"bearer" validate:"omitempty"`
	Key          string       `yaml:"key" validate:"omitempty"`
	Secret       string       `yaml:"secret" validate:"omitempty"`
	Certificates Certificates `yaml:"certificates" validate:"omitempty"`
}

type Certificates struct {
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
	CAFile   string `yaml:"CAFile"`
}
