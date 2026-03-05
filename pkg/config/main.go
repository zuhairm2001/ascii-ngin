package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds application configuration loaded from config.json.
type Config struct {
	GeminiAPIKey string `json:"GEMINI_API_KEY"`
}

var cfg Config
var loaded bool

// Load reads config.json from the current working directory and caches
// the result in package-level state. It returns an error if the file
// is missing, malformed, or if required fields are empty.
//
// TODO: resolve config path from ~/.ascii-ngin/config.json for distribution.
func Load() error {
	f, err := os.Open("config.json")
	if err != nil {
		return fmt.Errorf("config: open config.json: %w", err)
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
