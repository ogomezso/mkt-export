package config

type Config struct {
	Input  string `yaml:"inputpath" validate:"required"`
	Output string `yaml:"export" validate:"required"`
}
