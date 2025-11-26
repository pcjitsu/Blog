package config

import (
	// Used to convert Go structs to JSON string format and vice versa.
	"encoding/json"
	// Used for formatting strings and handling errors.
	"fmt"
	// Used for operating system functionality (reading/writing files, getting home dir).
	"os"
	// Used to handle file paths safely across different operating systems (Windows vs Mac/Linux).
	"path/filepath"
)

// configFileName is the name of the configuration file.
// The dot prefix (.) usually indicates a "hidden" file in Unix-like systems.
const configFileName = ".gatorconfig.json"

// Config represents the structure of the JSON file on disk.
// We use a struct to create a blueprint for the data.
type Config struct {
	// struct tags (the text in backticks) tell the JSON encoder/decoder
	// exactly which JSON key maps to which Go field.
	// e.g. "db_url" in JSON becomes DBURL in Go.
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// getConfigFilePath allows us to find the file regardless of who is logged in.
// It returns the full absolute path (e.g., /home/alice/.gatorconfig.json).
func getConfigFilePath() (string, error) {
	// os.UserHomeDir returns the home directory of the current user.
	// On Windows this might be C:\Users\Name, on Mac /Users/Name.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// filepath.Join is safer than manual string concatenation (homeDir + "/" + fileName).
	// It automatically adds the correct separator (\ for Windows, / for Linux).
	return filepath.Join(homeDir, configFileName), nil
}

// Read loads the file from disk and returns a Config struct.
func Read() (Config, error) {
	var cfg Config

	// 1. Get the path
	filePath, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	// 2. Read the raw bytes from the file.
	// os.ReadFile is a convenience function that opens, reads, and closes the file.
	data, err := os.ReadFile(filePath)
	if err != nil {
		// We wrap the error using %w to provide context (where did it fail?)
		// while preserving the original error type.
		return cfg, fmt.Errorf("failed to read config file from %s: %w", filePath, err)
	}

	// 3. Unmarshal (Decode) the JSON bytes into the Go struct.
	// We pass a pointer (&cfg) so Unmarshal can fill the empty struct with data.
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}

	return cfg, nil
}

// write is a private helper function (starts with lowercase 'w').
// It saves the current state of a Config struct to disk.
func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// MarshalIndent converts the Go struct into a JSON byte slice.
	// The arguments ("", "  ") tell it to pretty-print with 2-space indentation,
	// making the file human-readable.
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// Write the data to the file.
	// 0600 is a permission code (octal).
	// 6 (Read/Write) for the owner, 0 for group, 0 for others.
	// This ensures other users on the computer cannot read your config credentials.
	err = os.WriteFile(filePath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file to %s: %w", filePath, err)
	}

	return nil
}

// SetUser updates the current user and saves the file.
// IMPORTANT: This uses a "Pointer Receiver" (cfg *Config).
// If we used (cfg Config), we would be modifying a copy, not the original.
func (cfg *Config) SetUser(username string) error {
	// Update the field in memory
	cfg.CurrentUserName = username

	// Write the updated struct to disk immediately so the change persists.
	// We dereference (*cfg) because write expects a value, not a pointer.
	return write(*cfg)
}
