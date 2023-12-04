package armaria

// Manifest is the data structure that gets marshalled into the app manifest.
type Manifest struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Path              string   `json:"path"`
	HostType          string   `json:"type"`
	AllowedExtensions []string `json:"allowed_extensions"`
	AllowedOrigins    []string `json:"allowed_origins"`
}
