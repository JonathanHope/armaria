package config

import "errors"

// ErrConfigMissing is returned when the config file is missing.
var ErrConfigMissing = errors.New("config missing")
