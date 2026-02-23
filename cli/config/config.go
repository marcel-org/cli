package config

import (
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
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".marcel.yml")

	config := &Config{
		AuthToken:    "",
		WeekStartDay: "sunday",
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, nil
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return config, nil
	}

	if envToken := os.Getenv("MARCEL_TOKEN"); envToken != "" {
		config.AuthToken = envToken
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
