package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AdminAddress         string   `yaml:"admin_address"`
	DBFilename           string   `yaml:"db_filename"`
	ConfiguratorFilename string   `yaml:"configurator_filename"`
	ConfiguratorCommand  []string `yaml:"configurator_command"`
	IP                   *string  `yaml:"ip"`
	Hostname             *string  `yaml:"hostname"`
}

func New(filename string) (*Config, error) {
	cfg := &Config{}
	f, err := os.Open(filename)
	if err != nil {
		return cfg, fmt.Errorf(
			"failed to open configuration file %v, %w",
			filename, err,
		)
	}
	yd := yaml.NewDecoder(f)
	err = yd.Decode(cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to decode config, %w", err)
	}
	return cfg, nil
}
