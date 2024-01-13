package armariaapi

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jonathanhope/armaria/internal/paths"
	"github.com/jonathanhope/armaria/pkg/model"
)

// TODO: Windows needs registry entiries.

// InstallManifestFirefox installs the app manifest for Firefox.
func InstallManifestFirefox() error {
	path, err := paths.FirefoxManifest()
	if err != nil {
		return fmt.Errorf("error getting firefox manfiest path while installing manifest: %w", err)
	}

	return installManifest(path, false)
}

// InstallManifestChrome installs the app manifest for Firefox.
func InstallManifestChrome() error {
	path, err := paths.ChromeManifest()
	if err != nil {
		return fmt.Errorf("error getting chrome manfiest path while installing manifest: %w", err)
	}

	return installManifest(path, true)
}

// InstallManifestChromium installs the app manifest for Firefox.
func InstallManifestChromium() error {
	path, err := paths.ChromiumManifest()
	if err != nil {
		return fmt.Errorf("error getting chromium manfiest path while installing manifest: %w", err)
	}

	return installManifest(path, true)
}

// installManifest installs the app manifest.
func installManifest(path string, isChrome bool) error {
	hostPath, err := paths.Host()
	if err != nil {
		return fmt.Errorf("error getting host path while installing manifest: %w", err)
	}

	var buffer []byte
	if isChrome {
		manifest := armaria.ChromeManifest{
			Name:        "armaria",
			Description: "Armaria is a fast local-first bookmarks manager.",
			Path:        hostPath,
			HostType:    "stdio",
			AllowedOrigins: []string{
				"chrome-extension://cahkgigfdplmhgjbioakkgennhncioli/",
				"chrome-extension://fbnilfpngakppdkddndcnckolmlpghdf/",
			},
		}

		buffer, err = json.Marshal(manifest)
		if err != nil {
			return fmt.Errorf("error marshalling manifest while installing manifest: %w", err)
		}
	} else {
		manifest := armaria.FirefoxManifest{
			Name:              "armaria",
			Description:       "Armaria is a fast local-first bookmarks manager.",
			Path:              hostPath,
			HostType:          "stdio",
			AllowedExtensions: []string{"armaria@armaria.net"},
		}

		buffer, err = json.Marshal(manifest)
		if err != nil {
			return fmt.Errorf("error marshalling manifest while installing manifest: %w", err)
		}
	}

	handle, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating manifest file while installing manifest: %w", err)
	}
	defer handle.Close()

	_, err = handle.Write(buffer)
	if err != nil {
		return fmt.Errorf("error writing manfiest file contents while installing manifest: %w", err)
	}

	return nil
}
