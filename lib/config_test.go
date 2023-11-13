package lib

import (
	"os"
	"path"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	type test struct {
		goos          string
		configPath    string
		folderCreated bool
	}

	tests := []test{
		{
			goos:          "windows",
			configPath:    "~/AppData/Local/Armaria/armaria.toml",
			folderCreated: false,
		},
		{
			goos:          "linux",
			configPath:    "~/.armaria/armaria.toml",
			folderCreated: false,
		},
		{
			goos:          "darwin",
			configPath:    "~/Library/Application Support/Armaria/armaria.toml",
			folderCreated: false,
		},
	}

	userHome := func() (string, error) {
		return "~", nil
	}

	for _, tc := range tests {
		t.Run(tc.goos, func(t *testing.T) {
			folderCreated := false

			got, err := getConfigPath(tc.goos, userHome, path.Join)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			if folderCreated != tc.folderCreated {
				t.Fatalf("folder created: got %+v; want %+v", folderCreated, tc.folderCreated)
			}

			if got != tc.configPath {
				t.Errorf("db: got %+v; want %+v", got, tc.configPath)
			}
		})
	}
}

func TestGetDatabasePath(t *testing.T) {
	type test struct {
		inputPath     NullString
		configPath    string
		goos          string
		folder        string
		db            string
		folderCreated bool
	}

	tests := []test{
		{
			inputPath:     NullStringFromPtr(nil),
			configPath:    "",
			goos:          "windows",
			folder:        "~/AppData/Local/Armaria",
			db:            "~/AppData/Local/Armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			inputPath:     NullStringFromPtr(nil),
			configPath:    "",
			goos:          "linux",
			folder:        "~/.armaria",
			db:            "~/.armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			inputPath:     NullStringFromPtr(nil),
			configPath:    "",
			goos:          "darwin",
			folder:        "~/Library/Application Support/Armaria",
			db:            "~/Library/Application Support/Armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			inputPath:     NullStringFrom("bookmarks.db"),
			configPath:    "",
			db:            "bookmarks.db",
			folderCreated: false,
		},
	}

	userHome := func() (string, error) {
		return "~", nil
	}

	for _, tc := range tests {
		t.Run(tc.goos, func(t *testing.T) {
			folderCreated := false
			mkDirAll := func(path string, perm os.FileMode) error {
				folderCreated = true
				if path != tc.folder {
					t.Errorf("folder: got %+v; want %+v", path, tc.folder)
				}

				return nil
			}

			got, err := getDatabasePath(tc.inputPath, tc.configPath, tc.goos, mkDirAll, userHome, path.Join)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			if folderCreated != tc.folderCreated {
				t.Fatalf("folder created: got %+v; want %+v", folderCreated, tc.folderCreated)
			}

			if got != tc.db {
				t.Errorf("db: got %+v; want %+v", got, tc.db)
			}
		})
	}
}
