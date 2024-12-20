package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Address     string        `yaml:"address"`
		NonceTTL    time.Duration `yaml:"nonceTtl"`
		SecretKey   []byte        `yaml:"secretKey"`
		MaxDataSize int           `yaml:"maxDataSize"`
		MaxConn     int           `yaml:"maxConn"`
	} `yaml:"server"`

	Client struct {
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"client"`

	PoW struct {
		Difficulty int `yaml:"difficulty"`
	} `yaml:"pow"`

	Quotes struct {
		FilePath string `yaml:"filePath"`
	} `yaml:"quotes"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", "configs/config.yaml", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s: %w", "configs/config.yaml", err)
	}

	return &cfg, nil
}
