package lib

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

const configFilename = "armaria.toml"
const databaseFilename = "bookmarks.db"

// Config is the unmarshalled Config.
type Config struct {
	DB string `koanf:"db"`
}

// mkDirAllFn creates a directory if it doesn't already exist.
type mkDirAllFn func(path string, perm os.FileMode) error

// userHomeFn returns the home directory of the current user.
type userHomeFn func() (string, error)

// joinFn joins path segments together.
type joinFn func(elem ...string) string

// updateConfigFn is a callback to update the current config
type updateConfigFn func(config *Config)

// GetConfig gets the current config.
// If the sentinerl error ErrConfigMissing then it doesn't exist.
func GetConfig() (Config, error) {
	config := Config{}

	configPath, err := getConfigPath(runtime.GOOS, os.UserHomeDir, filepath.Join)
	if err != nil {
		return config, err
	}

	var k = koanf.New(".")
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return config, ErrConfigMissing
		} else {
			return config, err
		}
	}

	err = k.Unmarshal("", &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// UpdateConfig updates the current config.
// It will be created if it hasn't already been created.
func UpdateConfig(update updateConfigFn) error {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, ErrConfigMissing) {
		return err
	}

	if errors.Is(err, ErrConfigMissing) {
		folder, err := getFolderPath(runtime.GOOS, os.UserHomeDir, filepath.Join)
		if err != nil {
			return err
		}

		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	update(&config)

	var k = koanf.New(".")
	err = k.Load(structs.Provider(config, "koanf"), nil)
	if err != nil {
		return err
	}

	buffer, err := k.Marshal(toml.Parser())
	if err != nil {
		return err
	}

	configPath, err := getConfigPath(runtime.GOOS, os.UserHomeDir, filepath.Join)
	if err != nil {
		return err
	}

	handle, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer handle.Close()

	_, err = handle.Write(buffer)
	if err != nil {
		return err
	}

	return nil
}

// getDatabasePath gets the path to the bookmarks database.
// The path will be (in order of precedence):
// 1) The inputted path
// 2) The path in the config file
// 3) The default path (getFolderPath() + "bookmarks.db")
func getDatabasePath(inputPath NullString, configPath string, goos string, mkDirAll mkDirAllFn, userHome userHomeFn, join joinFn) (string, error) {
	if inputPath.Valid && inputPath.Dirty {
		return inputPath.String, nil
	} else if configPath != "" {
		return configPath, nil
	} else {
		folder, err := getFolderPath(goos, userHome, join)
		if err != nil {
			return "", err
		}

		if err = mkDirAll(folder, os.ModePerm); err != nil {
			return "", err
		}

		return join(folder, databaseFilename), nil
	}
}

// getConfigPath gets the path to the config file.
// The config file is a TOML file located at getFolderPath() + "bookmarks.db".
func getConfigPath(goos string, userHome userHomeFn, join joinFn) (string, error) {
	folder, err := getFolderPath(goos, userHome, join)
	if err != nil {
		return "", err
	}

	return join(folder, configFilename), nil
}

// getFolderPath gets the path to the folder the config (and by default) the database are stored.
// The folder is different per platform and maps to the following:
// - Linux: ~/.armaria
// - Windows: ~/AppData/Local/Armaria
// - Mac: ~/Library/Application Support/Armaria
func getFolderPath(goos string, userHome userHomeFn, join joinFn) (string, error) {
	home, err := userHome()
	if err != nil {
		return "", err
	}

	var folder string
	if goos == "linux" {
		folder = join(home, ".armaria")
	} else if goos == "windows" {
		folder = join(home, "AppData", "Local", "Armaria")
	} else if goos == "darwin" {
		folder = join(home, "Library", "Application Support", "Armaria")
	} else {
		panic("Unsupported operating system")
	}

	return folder, nil
}
