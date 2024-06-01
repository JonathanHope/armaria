package manifest

import (
	"encoding/json"
	"fmt"
	"os"
)

const name = "armaria"
const description = "Armaria is a fast local-first bookmarks manager."
const hostType = "stdio"
const FirefoxExtension = "armaria@armaria.net"
const ChromeExtension1 = "chrome-extension://cahkgigfdplmhgjbioakkgennhncioli/"
const ChromeExtension2 = "chrome-extension://fbnilfpngakppdkddndcnckolmlpghdf/"

// TODO: Windows needs registry entiries.

// InstallManifest installs the app manifest.
func InstallManifest(path string, hostPath string, manifestType ManifestType) error {
	var err error
	var buffer []byte

	if manifestType == ManifestChrome || manifestType == ManifestChromium {
		manifest := chromeManifest{
			Name:        name,
			Description: description,
			Path:        hostPath,
			HostType:    hostType,
			AllowedOrigins: []string{
				ChromeExtension1,
				ChromeExtension2,
			},
		}

		buffer, err = json.Marshal(manifest)
		if err != nil {
			return fmt.Errorf("error marshalling manifest while installing manifest: %w", err)
		}
	} else if manifestType == ManifestFirefox {
		manifest := firefoxManifest{
			Name:              name,
			Description:       description,
			Path:              hostPath,
			HostType:          hostType,
			AllowedExtensions: []string{FirefoxExtension},
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
