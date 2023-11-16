package paths

import (
	"os"
	"path"
	"testing"

	"github.com/jonathanhope/armaria/internal/null"
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

			got, err := configInternal(tc.goos, userHome, path.Join)
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
		inputPath     null.NullString
		configPath    string
		goos          string
		folder        string
		db            string
		folderCreated bool
	}

	tests := []test{
		{
			inputPath:     null.NullStringFromPtr(nil),
			configPath:    "",
			goos:          "windows",
			folder:        "~/AppData/Local/Armaria",
			db:            "~/AppData/Local/Armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			inputPath:     null.NullStringFromPtr(nil),
			configPath:    "",
			goos:          "linux",
			folder:        "~/.armaria",
			db:            "~/.armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			inputPath:     null.NullStringFromPtr(nil),
			configPath:    "",
			goos:          "darwin",
			folder:        "~/Library/Application Support/Armaria",
			db:            "~/Library/Application Support/Armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			inputPath:     null.NullStringFrom("bookmarks.db"),
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

			got, err := databaseInternal(tc.inputPath, tc.configPath, tc.goos, mkDirAll, userHome, path.Join)
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
