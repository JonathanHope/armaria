package armaria

import (
	"fmt"

	"github.com/jonathanhope/armaria/internal/manifest"
	"github.com/jonathanhope/armaria/internal/paths"
)

// InstallManifestFirefox installs the app manifest for Firefox.
func InstallManifestFirefox() error {
	path, err := paths.FirefoxManifest()
	if err != nil {
		return fmt.Errorf("error getting firefox manfiest path while installing manifest: %w", err)
	}
	hostPath, err := paths.Host()
	if err != nil {
		return fmt.Errorf("error getting host path while installing manifest: %w", err)
	}

	return manifest.InstallManifest(path, hostPath, manifest.ManifestFirefox)
}

// InstallManifestChrome installs the app manifest for Firefox.
func InstallManifestChrome() error {
	path, err := paths.ChromeManifest()
	if err != nil {
		return fmt.Errorf("error getting chrome manfiest path while installing manifest: %w", err)
	}
	hostPath, err := paths.Host()
	if err != nil {
		return fmt.Errorf("error getting host path while installing manifest: %w", err)
	}

	return manifest.InstallManifest(path, hostPath, manifest.ManifestChrome)
}

// InstallManifestChromium installs the app manifest for Firefox.
func InstallManifestChromium() error {
	path, err := paths.ChromiumManifest()
	if err != nil {
		return fmt.Errorf("error getting chromium manfiest path while installing manifest: %w", err)
	}
	hostPath, err := paths.Host()
	if err != nil {
		return fmt.Errorf("error getting host path while installing manifest: %w", err)
	}

	return manifest.InstallManifest(path, hostPath, manifest.ManifestChromium)
}
