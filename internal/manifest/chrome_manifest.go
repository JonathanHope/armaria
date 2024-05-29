package manifest

// chromeManifest is the data structure that gets marshalled into the app manifest for Chrome.
type chromeManifest struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Path           string   `json:"path"`
	HostType       string   `json:"type"`
	AllowedOrigins []string `json:"allowed_origins"`
}
