package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const APIEndpoint = "https://api.marcel.my"

type Config struct {
	AuthToken    string `yaml:"-"`
	WeekStartDay string `yaml:"week_start_day"`
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".marcel.yml")

	config := &Config{
		AuthToken:    "",
		WeekStartDay: "sunday",
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := config.Save(); err != nil {
			return nil, fmt.Errorf("failed to create config file at %s: %w", configPath, err)
		}
	} else {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file at %s: %w", configPath, err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("invalid YAML format in %s: %w", configPath, err)
		}
	}

	tokenPath := filepath.Join(homeDir, ".marcel.token")
	tokenBytes, err := os.ReadFile(tokenPath)
	if err == nil && len(tokenBytes) > 0 {
		config.AuthToken = strings.TrimSpace(string(tokenBytes))
	} else {
		config.AuthToken = os.Getenv("MARCEL_TOKEN")
	}

	if config.AuthToken == "" {
		return nil, fmt.Errorf("authentication token not found: neither ~/.marcel.token file nor MARCEL_TOKEN environment variable is set")
	}

	if config.WeekStartDay == "" {
		config.WeekStartDay = "sunday"
	}
	config.WeekStartDay = strings.ToLower(config.WeekStartDay)

	return config, nil
}

func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".marcel.yml")

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
