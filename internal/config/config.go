package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Assets []string `yaml:"assets"`
	Ignore []string `yaml:"ignore"`
}

func Load(projectPath string) (*Config, error) {
	configPath, err := findConfigPath(projectPath)
	if err != nil {
		return nil, err
	}
	if configPath == "" {
		return &Config{}, nil
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}
	return &cfg, nil
}

func findConfigPath(projectPath string) (string, error) {
	for _, name := range []string{"vole.yml", "vole.yaml"} {
		configPath := filepath.Join(projectPath, name)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("checking config file %s: %w", configPath, err)
		}
	}

	return "", nil
}
