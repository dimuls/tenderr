package operator

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type CORS struct {
	Enabled          bool   `yaml:"enabled"`
	AllowedOrigins   string `yaml:"allowedOrigins"`
	AllowedMethods   string `yaml:"allowedMethods"`
	AllowedHeaders   string `yaml:"allowedHeaders"`
	AllowCredentials bool   `yaml:"allowCredentials"`
}

type Config struct {
	Addr        string `yaml:"addr"`
	PostgresURL string `yaml:"postgresUrl"`
	CORS        CORS   `yaml:"cors"`
	LogLevel    string `yaml:"logLevel"`
}

func (c *Config) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config file: %w", err)
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(c)
	if err != nil {
		return fmt.Errorf("yaml deocde: %w", err)
	}

	return nil
}
