package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// GetConfig gets the current config.
// If the sentinel error ErrConfigMissing then it doesn't exist.
func GetConfig(configPath string) (Config, error) {
	config := Config{}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return config, nil
	}

	var k = koanf.New(".")
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return config, ErrConfigMissing
		} else {
			return config, fmt.Errorf("error loading config while getting config: %w", err)
		}
	}

	err := k.Unmarshal("", &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshalling config while getting config: %w", err)
	}

	return config, nil
}
