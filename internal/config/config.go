package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Assets []string `yaml:"assets`
	Ignore []string `yaml:"ignore"`
}

func Load(projectPath string) (*Config, error) {
	configPath := filepath.Join(projectPath, "vole.yml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w")
	}
	return &cfg, nil
}
