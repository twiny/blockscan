package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

// Config
type Config struct {
	// Rest
	Rest struct {
		Addr string `yaml:"address"`
	} `yaml:"rest"`

	// Indexer
	Indexer struct {
		Addr     string `yaml:"address"`
		Host     string `yaml:"host"`
		Token    string `yaml:"token"`
		Endpoint string `yaml:"endpoint"`
		Limiter  struct {
			Rate     int           `yaml:"rate"`
			Duration time.Duration `yaml:"duration"`
		} `yaml:"limiter"`
		Workers int           `yaml:"workers"`
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"indexer"`

	// Store
	Store struct {
		Path string `yaml:"path"`
	} `yaml:"store"`
}

// ParseConfig
func ParseConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

// Validate
func (c *Config) Validate() error {
	// TODO: validate config or use default values
	return nil
}
