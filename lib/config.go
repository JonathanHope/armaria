package lib

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const configFilename = "armaria.toml"
const databaseFilename = "bookmarks.db"

// mkDirAllFn creates a directory if it doesn't already exist.
type mkDirAllFn func(path string, perm os.FileMode) error

// userHomeFn returns the home directory of the current user.
type userHomeFn func() (string, error)

// joinFn joins path segments together.
type joinFn func(elem ...string) string

// GetDBPathConfig gets the path to the bookmarks database from the config file.
func GetDBPathConfig() (string, error) {
	config, err := getConfig(runtime.GOOS, os.UserHomeDir, filepath.Join)
	if err != nil && !errors.Is(err, ErrConfigMissing) {
		return "", err
	}

	if errors.Is(err, ErrConfigMissing) {
		return "", nil
	} else {
		return config.String("db"), nil
	}
}

// getConfig parses the config file.
// If the sentinerl error ErrConfigMissing then it doesn't exist.
func getConfig(goos string, userHome userHomeFn, join joinFn) (*koanf.Koanf, error) {
	configPath, err := getConfigPath(goos, userHome, join)
	if err != nil {
		return nil, err
	}

	var config = koanf.New(".")
	if err := config.Load(file.Provider(configPath), toml.Parser()); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil, ErrConfigMissing
		} else {
			return nil, err
		}
	}

	return config, nil
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
