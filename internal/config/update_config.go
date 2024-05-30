package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

// UpdateConfig updates the current config.
// It will be created if it hasn't already been created.
func UpdateConfig(folderPath string, configPath string, update UpdateConfigCallback) error {
	config, err := GetConfig(configPath)
	if err != nil && !errors.Is(err, ErrConfigMissing) {
		return fmt.Errorf("error getting config while updating config: %w", err)
	}

	if errors.Is(err, ErrConfigMissing) {
		err = os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error making config folder while updating config: %w", err)
		}
	}

	update(&config)

	var k = koanf.New(".")
	err = k.Load(structs.Provider(config, "koanf"), nil)
	if err != nil {
		return fmt.Errorf("error loading config while updating config: %w", err)
	}

	buffer, err := k.Marshal(toml.Parser())
	if err != nil {
		return fmt.Errorf("error marshalling config while updating config: %w", err)
	}

	handle, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("error creating config file while updating config: %w", err)
	}
	defer handle.Close()

	_, err = handle.Write(buffer)
	if err != nil {
		return fmt.Errorf("error writing config file contents while updating config: %w", err)
	}

	return nil
}
