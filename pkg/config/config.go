package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	// Game configuration
	Game struct {
		DataPath string `yaml:"dataPath"` // Path to game data files
		LogLevel string `yaml:"logLevel"`
	} `yaml:"game"`

	// Storage configuration
	Storage struct {
		Type      string `yaml:"type"`      // "file" for original file format
		PlayerDir string `yaml:"playerDir"` // Directory for player files
	} `yaml:"storage"`
}

// Load loads configuration from a YAML file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults if not specified
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 4000
	}
	if cfg.Game.LogLevel == "" {
		cfg.Game.LogLevel = "info"
	}
	if cfg.Storage.Type == "" {
		cfg.Storage.Type = "file"
	}

	return &cfg, nil
}
