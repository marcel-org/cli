package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIEndpoint string `yaml:"api_endpoint"`
	AuthToken   string `yaml:"auth_token"`
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".marcel.yml")

	config := &Config{
		APIEndpoint: "https://api.marcel.my",
		AuthToken:   "",
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

	if config.APIEndpoint == "" {
		config.APIEndpoint = "https://api.marcel.my"
	}

	if strings.HasPrefix(config.APIEndpoint, "~/") {
		config.APIEndpoint = filepath.Join(homeDir, config.APIEndpoint[2:])
	}

	if envToken := os.Getenv("MARCEL_TOKEN"); envToken != "" && config.AuthToken == "" {
		config.AuthToken = envToken
	}

	if envEndpoint := os.Getenv("MARCEL_API_ENDPOINT"); envEndpoint != "" {
		config.APIEndpoint = envEndpoint
	}

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
