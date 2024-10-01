package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"globalblackbox.io/gbx/models"

	"gopkg.in/yaml.v2"
)

// SaveConfig saves the configuration to ~/.gbx/config.yaml
func SaveConfig(config *models.Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to determine home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".gbx")
	configFile := filepath.Join(configDir, "config.yaml")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0700); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
