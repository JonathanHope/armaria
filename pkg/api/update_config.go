package armariaapi

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jonathanhope/armaria"
	"github.com/jonathanhope/armaria/internal/paths"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

// UpdateConfigCallback is a callback to update the current config
type UpdateConfigCallback func(config *armaria.Config)

// UpdateConfig updates the current config.
// It will be created if it hasn't already been created.
func UpdateConfig(update UpdateConfigCallback) error {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		return fmt.Errorf("error getting config while updating config: %w", err)
	}

	if errors.Is(err, armaria.ErrConfigMissing) {
		folder, err := paths.Folder(runtime.GOOS, os.UserHomeDir, filepath.Join)
		if err != nil {
			return fmt.Errorf("error getting config folder path while updating config: %w", err)
		}

		err = os.MkdirAll(folder, os.ModePerm)
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

	configPath, err := paths.Config()
	if err != nil {
		return fmt.Errorf("error getting config path while updating config: %w", err)
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
