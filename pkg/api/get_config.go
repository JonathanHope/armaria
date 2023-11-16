package armariaapi

import (
	"fmt"
	"strings"

	"github.com/jonathanhope/armaria"
	"github.com/jonathanhope/armaria/internal/paths"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// GetConfig gets the current config.
// If the sentinel error ErrConfigMissing then it doesn't exist.
func GetConfig() (armaria.Config, error) {
	config := armaria.Config{}

	configPath, err := paths.Config()
	if err != nil {
		return config, fmt.Errorf("error getting config path while getting config: %w", err)
	}

	var k = koanf.New(".")
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return config, armaria.ErrConfigMissing
		} else {
			return config, fmt.Errorf("error loading config while getting config: %w", err)
		}
	}

	err = k.Unmarshal("", &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshalling config while getting config: %w", err)
	}

	return config, nil
}
