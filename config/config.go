// Package config loads mux user configuration from the XDG config directory.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config contains user-selectable mux settings.
type Config struct {
	Theme       string                         `json:"theme"`
	Keybindings map[string]map[string][]string `json:"keybindings"`
}

// Path returns the XDG-compatible mux configuration path.
func Path() (string, error) {
	root := os.Getenv("XDG_CONFIG_HOME")
	if root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("find home directory: %w", err)
		}
		root = filepath.Join(home, ".config")
	}
	return filepath.Join(root, "mux", "config.json"), nil
}

// Load reads the user configuration. A missing file is treated as an empty
// configuration.
func Load() (Config, error) {
	path, err := Path()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("read config %s: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	var value Config
	if err := json.NewDecoder(file).Decode(&value); err != nil {
		return Config{}, fmt.Errorf("decode config %s: %w", path, err)
	}
	return value, nil
}
