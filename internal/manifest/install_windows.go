//go:build windows

package manifest

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// InstallManifest installs the app manifest.
func InstallManifest(path string, hostPath string, manifestType ManifestType) error {
	err := installManifest(path, hostPath, manifestType)
	if err != nil {
		return err
	}

	// Windows (unfortunately) requires registry entries here.
	if manifestType == ManifestChrome || manifestType == ManifestChromium {
		if err := writeKey(`Software\Google\Chrome\NativeMessagingHosts\armaria`, path); err != nil {
			return err
		}
	} else if manifestType == ManifestFirefox {
		if err := writeKey(`Software\Mozilla\NativeMessagingHosts\armaria`, path); err != nil {
			return err
		}
	}

	return nil
}

// writeKey writes value to registry key.
// It gets written to both local machine and current user.
func writeKey(path string, value string) error {
	localMachineKey, err := openKey(registry.LOCAL_MACHINE, path)
	if err != nil {
		return fmt.Errorf("error opening local machine key while installing manifest: %w", err)
	}
	defer localMachineKey.Close()

	if err := localMachineKey.SetStringValue("", value); err != nil {
		return fmt.Errorf("error writing local machine key while installing manifest: %w", err)
	}

	userKey, err := openKey(registry.CURRENT_USER, path)
	if err != nil {
		return fmt.Errorf("error opening current user key while installing manifest: %w", err)
	}
	defer userKey.Close()

	if err := userKey.SetStringValue("", value); err != nil {
		return fmt.Errorf("error writing current user key while installing manifest: %w", err)
	}

	return nil
}

// openKey opens a key for writing.
// It will create any parent keys necessary.
func openKey(root registry.Key, path string) (registry.Key, error) {
	tokens := strings.Split(path, `\`)
	keys := []registry.Key{root}
	for _, token := range tokens {
		key, _, err := registry.CreateKey(keys[len(keys)-1], token, registry.WRITE)
		if err != nil {
			return key, err
		}
		keys = append(keys, key)
	}

	for i, key := range keys {
		if i != 0 && i != len(keys)-1 {
			key.Close()
		}
	}

	return keys[len(keys)-1], nil
}
