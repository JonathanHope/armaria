package armaria

// ChromeManifest is the data structure that gets marshalled into the app manifest for Chrome.
type ChromeManifest struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Path           string   `json:"path"`
	HostType       string   `json:"type"`
	AllowedOrigins []string `json:"allowed_origins"`
}

// FirefoxManifest is the data structure that gets marshalled into the app manifest for Firefox.
type FirefoxManifest struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Path              string   `json:"path"`
	HostType          string   `json:"type"`
	AllowedExtensions []string `json:"allowed_extensions"`
}
