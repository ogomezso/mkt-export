package config

type Config struct {
	Input             string            `yaml:"inputpath" validate:"required"`
	Output            string            `yaml:"export" validate:"required"`
	MarketPlaceConfig MarketPlaceConfig `yaml:"marketplace" validate:"omitempty"`
}

type MarketPlaceConfig struct {
	AppKey       string      `yaml:"appkey" validate:"required"`
	MktPlaceUrl  string      `yaml:"url" validation:"required"`
	Credentials  Credentials `yaml:"credentials" validate:"omitempty"`
	SchemasPath  string      `yaml:"schemasuri" validate:"omitempty"`
	ExamplesPath string      `yaml:"examplesPath" validate:"omitempty"`
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
