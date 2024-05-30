package armaria

import (
	"fmt"

	"github.com/jonathanhope/armaria/internal/config"
	"github.com/jonathanhope/armaria/internal/paths"
)

type UpdateConfigCallback = config.UpdateConfigCallback

// UpdateConfig updates the current config.
// It will be created if it hasn't already been created.
func UpdateConfig(update UpdateConfigCallback) error {
	folderPath, err := paths.Folder()
	if err != nil {
		return fmt.Errorf("error getting config folder path while updating config: %w", err)
	}

	configPath, err := paths.Config()
	if err != nil {
		return fmt.Errorf("error getting config path while updating config: %w", err)
	}

	return config.UpdateConfig(folderPath, configPath, update)
}
