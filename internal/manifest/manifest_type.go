package manifest

// ManifestType is the type of manifest being installed.
type ManifestType int

const (
	ManifestFirefox ManifestType = iota + 1
	ManifestChrome
	ManifestChromium
)
