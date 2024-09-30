// cmd/config.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// SignupPlan represents the subscription plan details
type SignupPlan struct {
	Name   string `json:"name" yaml:"name"`
	Region string `json:"region,omitempty" yaml:"region,omitempty"`
}

// Config represents the structure of the configuration file
type Config struct {
	APIKey    string     `yaml:"api_key"`
	AccountID string     `yaml:"account_id"`
	Plan      SignupPlan `yaml:"plan"`
}

// SaveConfig saves the configuration to ~/.gbx/config.yaml
func SaveConfig(config *Config) error {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to determine home directory: %v", err)
	}

	// Define the config directory and file path
	configDir := filepath.Join(homeDir, ".gbx")
	configFile := filepath.Join(configDir, "config.yaml")

	// Create the config directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0700); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	// Marshal the config struct to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %v", err)
	}

	// Write the YAML data to the config file
	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
