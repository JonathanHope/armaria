package armaria

import (
	"fmt"

	"github.com/jonathanhope/armaria/internal/config"
	"github.com/jonathanhope/armaria/internal/paths"
)

// GetConfig gets the current config.
// If the sentinel error ErrConfigMissing then it doesn't exist.
func GetConfig() (Config, error) {
	configPath, err := paths.Config()
	if err != nil {
		return Config{}, fmt.Errorf("error getting config path while getting config: %w", err)
	}

	return config.GetConfig(configPath)
}
