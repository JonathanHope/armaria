// paths manages the paths to the various files Armaria cares about.
// Armaria stores its bookmarks in a SQLite DB and its config in a TOML file.
// Both of these files need to be stored somewhere.
// This file contains the logic to figure out where to store those files.
// It also keeps track of where manifest files need to be installed for browser extensions.
package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jonathanhope/armaria/internal/null"
)

// configFilename is the default name for the config file.
const configFilename = "armaria.toml"

// databaseFilename is the default name for the database.
const databaseFilename = "bookmarks.db"

// manfiestFilename is the default anmme for the app manifest.
const manifestFilename = "armaria.json"

// mkDirAllFn creates a directory if it doesn't already exist.
type mkDirAllFn func(path string, perm os.FileMode) error

// userHomeFn returns the home directory of the current user.
type userHomeFn func() (string, error)

// joinFn joins path segments together.
type joinFn func(elem ...string) string

// getenvFn gets an environment variable.
type getenvFn func(key string) string

// executableFn gets that path to the current executable.
type executableFn func() (string, error)

// dirFn gets the directory part of a path.
type dirFn func(string) string

// Database gets the path to the bookmarks database.
// The path will be (in order of precedence):
// 1) The inputted path
// 2) The path in the config file
// 3) The default path (getFolderPath() + "bookmarks.db")
func Database(inputPath null.NullString, configPath string) (string, error) {
	return databaseInternal(inputPath, configPath, runtime.GOOS, os.MkdirAll, os.UserHomeDir, filepath.Join, os.Getenv)
}

// databaseInternal allows DI for GetDatabasePath.
func databaseInternal(inputPath null.NullString, configPath string, goos string, mkDirAll mkDirAllFn, userHome userHomeFn, join joinFn, getenv getenvFn) (string, error) {
	if inputPath.Valid && inputPath.Dirty {
		return inputPath.String, nil
	} else if configPath != "" {
		return configPath, nil
	} else {
		folder, err := folderInternal(goos, userHome, join, getenv)
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
	return configInternal(runtime.GOOS, os.UserHomeDir, filepath.Join, os.Getenv, os.MkdirAll)
}

// configInternal allows DI for GetConfigPath.
func configInternal(goos string, userHome userHomeFn, join joinFn, getenv getenvFn, mkDirAll mkDirAllFn) (string, error) {
	folder, err := folderInternal(goos, userHome, join, getenv)
	if err != nil {
		return "", fmt.Errorf("error getting folder path while getting config path: %w", err)
	}

	if err = mkDirAll(folder, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating folder while getting config path: %w", err)
	}

	return join(folder, configFilename), nil
}

// Folder gets the path to the folder the config (and by default) the database are stored.
// The folder is different per platform and maps to the following:
// - Linux: ~/.armaria
// - Windows: ~/AppData/Local/Armaria
// - Mac: ~/Library/Application Support/Armaria
func Folder() (string, error) {
	return folderInternal(runtime.GOOS, os.UserHomeDir, filepath.Join, os.Getenv)
}

// folderInternal allows DI for Folder.
func folderInternal(goos string, userHome userHomeFn, join joinFn, getenv getenvFn) (string, error) {
	home, err := snapOrHome(userHome, getenv)
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

// snapOrHome returns the user common directory if running under snap.
// Otherise it returns the usual home directory.
// We need to make sure to use the common directory under snap so we don't lose data between revisions.
func snapOrHome(userHome userHomeFn, getenv getenvFn) (string, error) {
	snapDir := getenv("SNAP_USER_COMMON")
	if snapDir != "" {
		return snapDir, nil
	}

	return userHome()
}

// Host will get the path to the native messaging host.
// The extensions need an absolute path to it in order to work.
func Host() (string, error) {
	return hostInternal(runtime.GOOS, os.Getenv, os.Executable, filepath.Dir, filepath.Join)
}

func hostInternal(goos string, getenv getenvFn, executable executableFn, dir dirFn, join joinFn) (string, error) {
	// The snap path needs to be a hard coded special case.
	// There doesn't appear to be way to intuit it.
	// Additionally snap will namespace the host with "armaria".

	snapHome := getenv("SNAP_REAL_HOME")
	if snapHome != "" {
		return "/snap/bin/armaria", nil
	}

	// Otherwise we can assume the host is installed alongside the CLI.

	ex, err := executable()
	if err != nil {
		return "", fmt.Errorf("error getting current executable while getting host path: %w", err)
	}

	hostExe := "armaria"
	if goos == "windows" {
		hostExe = "armaria.exe"
	}

	return join(dir(ex), hostExe), nil
}

// FirefoxManifest gets the path to the Firefox app manifest.
// The path is different per platform and maps to the following:
// - Linux: ~/.mozilla/native-messaging-hosts
// - Windows: ~/AppData/Local/Armaria
// - Mac: ~/Library/Application Support/Mozilla/NativeMessagingHosts
func FirefoxManifest() (string, error) {
	return firefoxManifestInternal(runtime.GOOS, os.Getenv, os.UserHomeDir, filepath.Join, os.MkdirAll)
}

// firefoxManifestInternal allows DI for FirefoxManifest.
func firefoxManifestInternal(goos string, getenv getenvFn, userHome userHomeFn, join joinFn, mkDirAll mkDirAllFn) (string, error) {
	home, err := realHome(getenv, userHome)
	if err != nil {
		return "", fmt.Errorf("error getting real home dir while getting firefox manifest path: %w", err)
	}

	var folder string
	if goos == "linux" {
		folder = join(home, ".mozilla", "native-messaging-hosts")
	} else if goos == "windows" {
		// The manifest can be anywhere in Windows, but it needs a supporting registry entry.
		folder = join(home, "AppData", "Local", "Armaria")
	} else if goos == "darwin" {
		folder = join(home, "Library", "Application Support", "Mozilla", "NativeMessagingHosts")
	} else {
		panic("Unsupported operating system")
	}

	if err = mkDirAll(folder, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating folder while getting firefox manifest path: %w", err)
	}

	return join(folder, manifestFilename), nil
}

// ChromeManifest gets the path to the Chrome app manifest.
// The path is different per platform and maps to the following:
// - Linux: ~/.config/google-chrome/NativeMessagingHosts
// - Windows: ~/AppData/Local/Armaria
// - Mac: ~/Library/Application Support/Google/Chrome
func ChromeManifest() (string, error) {
	return chromeManifestInternal(runtime.GOOS, os.Getenv, os.UserHomeDir, filepath.Join, os.MkdirAll)
}

// chromeManifestInternal allows DI for ChromeManifest.
func chromeManifestInternal(goos string, getenv getenvFn, userHome userHomeFn, join joinFn, mkDirAll mkDirAllFn) (string, error) {
	home, err := realHome(getenv, userHome)
	if err != nil {
		return "", fmt.Errorf("error getting real home dir while getting chrome manifest path: %w", err)
	}

	var folder string
	if goos == "linux" {
		folder = join(home, ".config", "google-chrome", "NativeMessagingHosts")
	} else if goos == "windows" {
		// The manifest can be anywhere in Windows, but it needs a supporting registry entry.
		folder = join(home, "AppData", "Local", "Armaria")
	} else if goos == "darwin" {
		folder = join(home, "Library", "Application Support", "Google", "Chrome", "NativeMessagingHosts")
	} else {
		panic("Unsupported operating system")
	}

	if err = mkDirAll(folder, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating folder while getting chrome manifest path: %w", err)
	}

	return join(folder, manifestFilename), nil
}

// Chromium Manifest gets the path to the Chromium app manifest.
// The path is different per platform and maps to the following:
// - Linux: ~/.config/chromium/NativeMessagingHosts
// - Windows: ~/AppData/Local/Armaria
// - Mac: ~/Library/Application Support/Chromium
func ChromiumManifest() (string, error) {
	return chromiumManifestInternal(runtime.GOOS, os.Getenv, os.UserHomeDir, filepath.Join, os.MkdirAll)
}

// chromiumManifestInternal allows DI for ChromiumManifest.
func chromiumManifestInternal(goos string, getenv getenvFn, userHome userHomeFn, join joinFn, mkDirAll mkDirAllFn) (string, error) {
	home, err := realHome(getenv, userHome)
	if err != nil {
		return "", fmt.Errorf("error getting real home dir while getting chromium manifest path: %w", err)
	}

	var folder string
	if goos == "linux" {
		folder = join(home, ".config", "chromium", "NativeMessagingHosts")
	} else if goos == "windows" {
		// The manifest can be anywhere in Windows, but it needs a supporting registry entry.
		folder = join(home, "AppData", "Local", "Armaria")
	} else if goos == "darwin" {
		folder = join(home, "Library", "Application Support", "Chromium", "NativeMessagingHosts")
	} else {
		panic("Unsupported operating system")
	}

	if err = mkDirAll(folder, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating folder while getting chromium manifest path: %w", err)
	}

	return join(folder, manifestFilename), nil
}

// realHome returns the true home directory of the current user.
// Snap will replace the $HOME env var with a sandboxed directory.
func realHome(getenv getenvFn, userHome userHomeFn) (string, error) {
	home := getenv("SNAP_REAL_HOME")

	if home == "" {
		home, err := userHome()
		if err != nil {
			return "", fmt.Errorf("error getting home dir: %w", err)
		}
		return home, nil
	}

	return home, nil
}
