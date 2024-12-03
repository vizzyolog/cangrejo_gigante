package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Address   string        `yaml:"address"`
		NonceTTL  time.Duration `yaml:"nonceTTL"`
		SecretKey string        `yaml:"secretKey"`
	} `yaml:"server"`

	PoW struct {
		Difficulty int `yaml:"difficulty"`
	} `yaml:"pow"`

	Quotes struct {
		FilePath string `yaml:"file_path"`
	} `yaml:"quotes"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
