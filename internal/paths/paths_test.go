package paths

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/jonathanhope/armaria/internal/null"
)

func TestGetConfigPath(t *testing.T) {
	type test struct {
		goos          string
		configPath    string
		snapCommonDir string
	}

	tests := []test{
		{
			goos:       "windows",
			configPath: "~/AppData/Local/Armaria/armaria.toml",
		},
		{
			goos:       "linux",
			configPath: "~/.armaria/armaria.toml",
		},
		{
			goos:          "linux",
			configPath:    "~/snap/.armaria/armaria.toml",
			snapCommonDir: "~/snap",
		},
		{
			goos:       "darwin",
			configPath: "~/Library/Application Support/Armaria/armaria.toml",
		},
	}

	userHome := func() (string, error) {
		return "~", nil
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("GOOS: %s, SNAP_USER_COMMON: %s", tc.goos, tc.snapCommonDir), func(t *testing.T) {
			getenv := func(key string) string {
				return tc.snapCommonDir
			}

			got, err := configInternal(tc.goos, userHome, path.Join, getenv)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			if got != tc.configPath {
				t.Errorf("configPath: got %+v; want %+v", got, tc.configPath)
			}
		})
	}
}

func TestGetDatabasePath(t *testing.T) {
	type test struct {
		inputPath     null.NullString
		configPath    string
		goos          string
		folderPath    string
		dbPath        string
		folderCreated bool
		snapCommonDir string
	}

	tests := []test{
		{
			inputPath:     null.NullStringFromPtr(nil),
			configPath:    "",
			goos:          "windows",
			folderPath:    "~/AppData/Local/Armaria",
			dbPath:        "~/AppData/Local/Armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			inputPath:     null.NullStringFromPtr(nil),
			configPath:    "",
			goos:          "linux",
			folderPath:    "~/.armaria",
			dbPath:        "~/.armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			inputPath:     null.NullStringFromPtr(nil),
			configPath:    "",
			goos:          "linux",
			folderPath:    "~/snap/.armaria",
			dbPath:        "~/snap/.armaria/bookmarks.db",
			folderCreated: true,
			snapCommonDir: "~/snap",
		},
		{
			inputPath:     null.NullStringFromPtr(nil),
			configPath:    "",
			goos:          "darwin",
			folderPath:    "~/Library/Application Support/Armaria",
			dbPath:        "~/Library/Application Support/Armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			inputPath:     null.NullStringFrom("bookmarks.db"),
			configPath:    "",
			dbPath:        "bookmarks.db",
			folderCreated: false,
		},
	}

	userHome := func() (string, error) {
		return "~", nil
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("GOOS: %s, SNAP_USER_COMMON: %s", tc.goos, tc.snapCommonDir), func(t *testing.T) {
			folderCreated := false

			mkDirAll := func(path string, perm os.FileMode) error {
				folderCreated = true
				if path != tc.folderPath {
					t.Errorf("folder: got %+v; want %+v", path, tc.folderPath)
				}

				return nil
			}

			getenv := func(key string) string {
				return tc.snapCommonDir
			}

			got, err := databaseInternal(tc.inputPath, tc.configPath, tc.goos, mkDirAll, userHome, path.Join, getenv)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			if folderCreated != tc.folderCreated {
				t.Fatalf("folder created: got %+v; want %+v", folderCreated, tc.folderCreated)
			}

			if got != tc.dbPath {
				t.Errorf("dbPath: got %+v; want %+v", got, tc.dbPath)
			}
		})
	}
}

func TestGetFolderPath(t *testing.T) {
	type test struct {
		goos          string
		folderPath    string
		snapCommonDir string
	}

	tests := []test{
		{
			goos:       "windows",
			folderPath: "~/AppData/Local/Armaria",
		},
		{
			goos:       "linux",
			folderPath: "~/.armaria",
		},
		{
			goos:          "linux",
			folderPath:    "~/snap/.armaria",
			snapCommonDir: "~/snap",
		},
		{
			goos:       "darwin",
			folderPath: "~/Library/Application Support/Armaria",
		},
	}

	userHome := func() (string, error) {
		return "~", nil
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("GOOS: %s, SNAP_USER_COMMON: %s", tc.goos, tc.snapCommonDir), func(t *testing.T) {
			getenv := func(key string) string {
				return tc.snapCommonDir
			}

			got, err := folderInternal(tc.goos, userHome, path.Join, getenv)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			if got != tc.folderPath {
				t.Errorf("db: got %+v; want %+v", got, tc.folderPath)
			}
		})
	}
}

func TestHostPath(t *testing.T) {
	type test struct {
		snapRealHome string
		hostPath     string
	}

	tests := []test{
		{
			snapRealHome: "",
			hostPath:     "/usr/bin/armaria-host",
		},
		{
			snapRealHome: "/snap",
			hostPath:     "/snap/bin/armaria.armaria-host",
		},
	}

	executable := func() (string, error) {
		return "/usr/bin/armaria", nil
	}

	dir := func(path string) string {
		if path == "/usr/bin/armaria" {
			return "/usr/bin"
		}

		return ""
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("SNAP_REAL_HOME: %s", tc.snapRealHome), func(t *testing.T) {
			getenv := func(key string) string {
				return tc.snapRealHome
			}

			got, err := hostInternal(getenv, executable, dir, path.Join)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			if got != tc.hostPath {
				t.Errorf("hostPath: got %+v; want %+v", got, tc.hostPath)
			}
		})
	}
}

func TestFirefoxManifestPath(t *testing.T) {
	type test struct {
		goos          string
		folderPath    string
		folderCreated bool
		snapRealHome  string
		manifestPath  string
	}

	tests := []test{
		{
			goos:          "windows",
			folderPath:    "~/AppData/Local/Armaria",
			folderCreated: true,
			manifestPath:  "~/AppData/Local/Armaria/armaria.json",
		},
		{
			goos:          "linux",
			folderPath:    "~/.mozilla/native-messaging-hosts",
			folderCreated: true,
			manifestPath:  "~/.mozilla/native-messaging-hosts/armaria.json",
		},
		{
			goos:          "linux",
			folderPath:    "~/snap/.mozilla/native-messaging-hosts",
			folderCreated: true,
			snapRealHome:  "~/snap",
			manifestPath:  "~/snap/.mozilla/native-messaging-hosts/armaria.json",
		},
		{
			goos:          "darwin",
			folderPath:    "~/Library/Application Support/Mozilla/NativeMessagingHosts",
			folderCreated: true,
			manifestPath:  "~/Library/Application Support/Mozilla/NativeMessagingHosts/armaria.json",
		},
	}

	userHome := func() (string, error) {
		return "~", nil
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("GOOS: %s, SNAP_REAL_HOME: %s", tc.goos, tc.snapRealHome), func(t *testing.T) {
			folderCreated := false

			mkDirAll := func(path string, perm os.FileMode) error {
				folderCreated = true
				if path != tc.folderPath {
					t.Errorf("folder: got %+v; want %+v", path, tc.folderPath)
				}

				return nil
			}

			getenv := func(key string) string {
				return tc.snapRealHome
			}

			got, err := firefoxManifestInternal(tc.goos, getenv, userHome, path.Join, mkDirAll)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			if folderCreated != tc.folderCreated {
				t.Fatalf("folder created: got %+v; want %+v", folderCreated, tc.folderCreated)
			}

			if got != tc.manifestPath {
				t.Errorf("manfiestPath: got %+v; want %+v", got, tc.manifestPath)
			}
		})
	}
}

func TestChromeManifestPath(t *testing.T) {
	type test struct {
		goos          string
		folderPath    string
		folderCreated bool
		snapRealHome  string
		manifestPath  string
	}

	tests := []test{
		{
			goos:          "windows",
			folderPath:    "~/AppData/Local/Armaria",
			folderCreated: true,
			manifestPath:  "~/AppData/Local/Armaria/armaria.json",
		},
		{
			goos:          "linux",
			folderPath:    "~/.config/google-chrome/NativeMessagingHosts",
			folderCreated: true,
			manifestPath:  "~/.config/google-chrome/NativeMessagingHosts/armaria.json",
		},
		{
			goos:          "linux",
			folderPath:    "~/snap/.config/google-chrome/NativeMessagingHosts",
			folderCreated: true,
			snapRealHome:  "~/snap",
			manifestPath:  "~/snap/.config/google-chrome/NativeMessagingHosts/armaria.json",
		},
		{
			goos:          "darwin",
			folderPath:    "~/Library/Application Support/Google/Chrome/NativeMessagingHosts",
			folderCreated: true,
			manifestPath:  "~/Library/Application Support/Google/Chrome/NativeMessagingHosts/armaria.json",
		},
	}

	userHome := func() (string, error) {
		return "~", nil
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("GOOS: %s, SNAP_REAL_HOME: %s", tc.goos, tc.snapRealHome), func(t *testing.T) {
			folderCreated := false

			mkDirAll := func(path string, perm os.FileMode) error {
				folderCreated = true
				if path != tc.folderPath {
					t.Errorf("folder: got %+v; want %+v", path, tc.folderPath)
				}

				return nil
			}

			getenv := func(key string) string {
				return tc.snapRealHome
			}

			got, err := chromeManifestInternal(tc.goos, getenv, userHome, path.Join, mkDirAll)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			if folderCreated != tc.folderCreated {
				t.Fatalf("folder created: got %+v; want %+v", folderCreated, tc.folderCreated)
			}

			if got != tc.manifestPath {
				t.Errorf("manfiestPath: got %+v; want %+v", got, tc.manifestPath)
			}
		})
	}
}
