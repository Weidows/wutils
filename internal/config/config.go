// Package config provides configuration management for wutils.
//
// It handles loading, saving, merging, and hot-reloading of YAML configuration files.
package config

// AppConfig holds application-level settings.
type AppConfig struct {
	Name    string `default:"wutils"`
	Version string `default:"1.0.0"`
	Debug   bool   `default:"false"`
}

// LoggingConfig holds logging settings.
type LoggingConfig struct {
	Level  string `default:"info"`
	Format string `default:"json"`
}

// DSGConfig holds Disk Sleep Guard settings.
type DSGConfig struct {
	Parallel bool     `default:"true"`
	Disk     []string `required:"true"`
	Delay    int      `default:"30"`
}

// OLPattern defines a window matching rule for the Opacity Listener.
type OLPattern struct {
	Title   string
	Opacity byte
}

// OLConfig holds Opacity Listener settings.
type OLConfig struct {
	Parallel bool        `default:"true"`
	Delay    int         `default:"2"`
	Patterns []OLPattern
}

// BufferConfig holds IO buffer filesystem settings.
type BufferConfig struct {
	Enable        bool   `default:"false"`
	MemoryLimit   int64  `default:"67108864"`
	FlushInterval int    `default:"10"`
	Strategy      string `default:"balanced"`
}

// CmdConfig groups all command-specific configurations.
type CmdConfig struct {
	Dsg    DSGConfig
	Ol     OLConfig
	Buffer BufferConfig
}

// Config is the top-level configuration structure for wutils.
type Config struct {
	App     AppConfig
	Logging LoggingConfig
	Refresh int `default:"10"`
	Cmd     CmdConfig
}
