package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds application configuration loaded from config.json.
type Config struct {
	GeminiAPIKey string `json:"GEMINI_API_KEY"`
}

var cfg Config
var loaded bool

// Load finds config.json by walking up from the current working directory
// until it reaches the filesystem root, then reads and caches the result.
// It returns an error if the file is not found, malformed, or if required
// fields are empty.
//
// TODO: resolve config path from ~/.ascii-ngin/config.json for distribution.
func Load() error {
	path, err := findConfig("config.json")
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("config: open %s: %w", path, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return fmt.Errorf("config: decode config.json: %w", err)
	}

	if cfg.GeminiAPIKey == "" {
		return fmt.Errorf("config: GEMINI_API_KEY is empty")
	}

	loaded = true
	return nil
}

// Get returns the cached Config. It panics if Load has not been called
// successfully, ensuring misuse is caught early during development.
func Get() Config {
	if !loaded {
		panic("config: Get() called before Load()")
	}
	return cfg
}

// findConfig walks up from the current working directory looking for
// a file with the given name. Returns the absolute path if found.
func findConfig(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %w", err)
	}

	for {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("%s not found in any parent directory", name)
		}
		dir = parent
	}
}
