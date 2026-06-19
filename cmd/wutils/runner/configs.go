package runner

import (
	"github.com/Weidows/wutils/internal/config"
)

// ConfigTemplate returns the default config template YAML bytes from the executable's config/ directory.
func ConfigTemplate() ([]byte, error) {
	return config.Template()
}

// Re-export config types for backward compatibility during migration.
// These will be removed once all code moves to internal/config.
type (
	AppConfig     = config.AppConfig
	LoggingConfig = config.LoggingConfig
	BufferConfig  = config.BufferConfig
	Config        = config.Config
)
