package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// configFileName is the hardcoded name of the config file.
const configFileName = ".gatorconfig.json"

// Config corresponds to the JSON structure of the config file.
type Config struct {
	DBURL           string `json:"db_url"`            // Maps to "db_url" in JSON
	CurrentUserName string `json:"current_user_name"` // Maps to "current_user_name" in JSON
}

// getConfigFilePath constructs the full path: /home/user/.gatorconfig.json
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, configFileName), nil
}

// Read loads the configuration from disk into memory.
func Read() (Config, error) {
	var cfg Config

	filePath, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	// Read the file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file from %s: %w", filePath, err)
	}

	// Parse JSON into the Go struct
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}

	return cfg, nil
}

// write saves the current Config struct to disk.
func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Convert struct back to JSON with indentation for readability
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// Write file with restricted permissions (0600 = read/write for owner only)
	err = os.WriteFile(filePath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file to %s: %w", filePath, err)
	}

	return nil
}

// SetUser updates the CurrentUserName field and immediately saves to disk.
func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)
}
