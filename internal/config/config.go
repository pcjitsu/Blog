package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// configFileName is the name of the configuration file in the user's home directory.
const configFileName = ".gatorconfig.json"

// Config represents the structure of the ~/.gatorconfig.json file.
// Struct tags are used to map Go fields to JSON keys.
type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// getConfigFilePath determines the full, absolute path to the configuration file
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, configFileName), nil
}

func Read() (Config, error) {

	var cfg Config
	filePath, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}
	// Read the file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file from %s: %w", filePath, err)
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}
	return cfg, nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// Write the JSON data to the file, creating it if it doesn't exist
	// and overwriting existing contents (0600 is owner read/write permissions)
	err = os.WriteFile(filePath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file to %s: %w", filePath, err)
	}

	return nil
}

// Pointer Receiver to mutuate the struct
func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)
}
