//go:build !windows

package manifest

// InstallManifest installs the app manifest.
func InstallManifest(path string, hostPath string, manifestType ManifestType) error {
	return installManifest(path, hostPath, manifestType)
}
