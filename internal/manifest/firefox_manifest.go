package manifest

// firefoxManifest is the data structure that gets marshalled into the app manifest for Firefox.
type firefoxManifest struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Path              string   `json:"path"`
	HostType          string   `json:"type"`
	AllowedExtensions []string `json:"allowed_extensions"`
}
