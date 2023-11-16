package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jonathanhope/armaria/internal/null"
)

// Armaria stores its bookmarks in a SQLite DB and its config in a TOML file.
// Both of these files need to be stored somewhere.
// This file contains the logic to figure out where to store those files.

// configFilename is the default name for the config file.
const configFilename = "armaria.toml"

// databaseFilename is the default name for the database.
const databaseFilename = "bookmarks.db"

// mkDirAllFn creates a directory if it doesn't already exist.
type mkDirAllFn func(path string, perm os.FileMode) error

// userHomeFn returns the home directory of the current user.
type userHomeFn func() (string, error)

// joinFn joins path segments together.
type joinFn func(elem ...string) string

// Database gets the path to the bookmarks database.
// The path will be (in order of precedence):
// 1) The inputted path
// 2) The path in the config file
// 3) The default path (getFolderPath() + "bookmarks.db")
func Database(inputPath null.NullString, configPath string) (string, error) {
	return databaseInternal(inputPath, configPath, runtime.GOOS, os.MkdirAll, os.UserHomeDir, filepath.Join)
}

// databaseInternal allows DI for GetDatabasePath.
func databaseInternal(inputPath null.NullString, configPath string, goos string, mkDirAll mkDirAllFn, userHome userHomeFn, join joinFn) (string, error) {
	if inputPath.Valid && inputPath.Dirty {
		return inputPath.String, nil
	} else if configPath != "" {
		return configPath, nil
	} else {
		folder, err := Folder(goos, userHome, join)
		if err != nil {
			return "", fmt.Errorf("error getting folder path while getting database path: %w", err)
		}

		if err = mkDirAll(folder, os.ModePerm); err != nil {
			return "", fmt.Errorf("error creating folder while getting database path: %w", err)
		}

		return join(folder, databaseFilename), nil
	}
}

// Config gets the path to the config file.
// The config file is a TOML file located at getFolderPath() + "bookmarks.db".
func Config() (string, error) {
	return configInternal(runtime.GOOS, os.UserHomeDir, filepath.Join)
}

// configInternal allows DI for GetConfigPath.
func configInternal(goos string, userHome userHomeFn, join joinFn) (string, error) {
	folder, err := Folder(goos, userHome, join)
	if err != nil {
		return "", fmt.Errorf("error getting folder path while getting config path: %w", err)
	}

	return join(folder, configFilename), nil
}

// Folder gets the path to the folder the config (and by default) the database are stored.
// The folder is different per platform and maps to the following:
// - Linux: ~/.armaria
// - Windows: ~/AppData/Local/Armaria
// - Mac: ~/Library/Application Support/Armaria
func Folder(goos string, userHome userHomeFn, join joinFn) (string, error) {
	home, err := userHome()
	if err != nil {
		return "", fmt.Errorf("error getting home path while getting folder path: %w", err)
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
