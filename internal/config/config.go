// Package config provides configuration management for wutils.
//
// Architecture
//
// Config has three layers, each overriding the previous:.
//  1. Compiled-in defaults (DefaultConfig)
//  2. User config file (~/.config/wutils/app.yml)
//  3. Runtime overrides (CLI flags / GUI inputs)
//
// The Watcher provides hot-reload via polling and implements ConfigProvider
// so that long-running services can pick up config changes at runtime.

package config

import "fmt"

// AppConfig holds application-level settings.
type AppConfig struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
	Debug   bool   `yaml:"debug" json:"debug"`
}

// LoggingConfig holds logging settings.
type LoggingConfig struct {
	Level  string `yaml:"level" json:"level"`
	Format string `yaml:"format" json:"format"`
}

// DSGConfig holds Disk Sleep Guard settings.
type DSGConfig struct {
	Parallel bool     `yaml:"parallel" json:"parallel"`
	Disk     []string `yaml:"disk" json:"disk"`
	Delay    int      `yaml:"delay" json:"delay"`
}

// OLPattern defines a window matching rule for the Opacity Listener.
type OLPattern struct {
	Title   string `yaml:"title" json:"title"`
	Opacity byte   `yaml:"opacity" json:"opacity"`
}

// OLConfig holds Opacity Listener settings.
type OLConfig struct {
	Parallel bool        `yaml:"parallel" json:"parallel"`
	Delay    int         `yaml:"delay" json:"delay"`
	Patterns []OLPattern `yaml:"patterns" json:"patterns"`
}

// BufferConfig holds IO buffer filesystem settings.
type BufferConfig struct {
	Enable        bool   `yaml:"enable" json:"enable"`
	MemoryLimit   int64  `yaml:"memory_limit" json:"memoryLimit"`
	FlushInterval int    `yaml:"flush_interval" json:"flushInterval"`
	Strategy      string `yaml:"strategy" json:"strategy"`
}

// CmdConfig groups all command-specific configurations.
type CmdConfig struct {
	Dsg    DSGConfig    `yaml:"dsg" json:"dsg"`
	Ol     OLConfig     `yaml:"ol" json:"ol"`
	Buffer BufferConfig `yaml:"buffer" json:"buffer"`
}

// Config is the top-level configuration structure for wutils.
type Config struct {
	App     AppConfig     `yaml:"app" json:"app"`
	Logging LoggingConfig `yaml:"logging" json:"logging"`
	Refresh int           `yaml:"refresh" json:"refresh"`
	Cmd     CmdConfig     `yaml:"cmd" json:"cmd"`
}

// Validate checks the configuration for common errors.
func (c *Config) Validate() error {
	if c.Refresh < 1 {
		return fmt.Errorf("refresh interval must be >= 1 second, got %d", c.Refresh)
	}
	if len(c.Cmd.Dsg.Disk) == 0 {
		return fmt.Errorf("dsg: at least one disk must be configured")
	}
	if c.Cmd.Dsg.Delay < 1 {
		return fmt.Errorf("dsg: delay must be >= 1 second, got %d", c.Cmd.Dsg.Delay)
	}
	if c.Cmd.Ol.Delay < 1 {
		return fmt.Errorf("ol: delay must be >= 1 second, got %d", c.Cmd.Ol.Delay)
	}
	switch c.Cmd.Buffer.Strategy {
	case "monitoring", "defrag", "download", "migration", "balanced":
		// valid
	case "":
		c.Cmd.Buffer.Strategy = "balanced"
	default:
		return fmt.Errorf("buffer: unknown strategy %q", c.Cmd.Buffer.Strategy)
	}
	return nil
}

// Merge copies non-zero fields from other into c.
// other's values take precedence where set.
func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}
	if other.App.Name != "" {
		c.App.Name = other.App.Name
	}
	if other.App.Version != "" {
		c.App.Version = other.App.Version
	}
	c.App.Debug = c.App.Debug || other.App.Debug

	if other.Logging.Level != "" {
		c.Logging.Level = other.Logging.Level
	}
	if other.Logging.Format != "" {
		c.Logging.Format = other.Logging.Format
	}

	if other.Refresh > 0 {
		c.Refresh = other.Refresh
	}

	if other.Cmd.Dsg.Parallel {
		c.Cmd.Dsg.Parallel = true
	}
	if len(other.Cmd.Dsg.Disk) > 0 {
		c.Cmd.Dsg.Disk = other.Cmd.Dsg.Disk
	}
	if other.Cmd.Dsg.Delay > 0 {
		c.Cmd.Dsg.Delay = other.Cmd.Dsg.Delay
	}

	if other.Cmd.Ol.Parallel {
		c.Cmd.Ol.Parallel = true
	}
	if other.Cmd.Ol.Delay > 0 {
		c.Cmd.Ol.Delay = other.Cmd.Ol.Delay
	}
	if len(other.Cmd.Ol.Patterns) > 0 {
		c.Cmd.Ol.Patterns = other.Cmd.Ol.Patterns
	}

	c.Cmd.Buffer.Enable = c.Cmd.Buffer.Enable || other.Cmd.Buffer.Enable
	if other.Cmd.Buffer.MemoryLimit > 0 {
		c.Cmd.Buffer.MemoryLimit = other.Cmd.Buffer.MemoryLimit
	}
	if other.Cmd.Buffer.FlushInterval > 0 {
		c.Cmd.Buffer.FlushInterval = other.Cmd.Buffer.FlushInterval
	}
	if other.Cmd.Buffer.Strategy != "" {
		c.Cmd.Buffer.Strategy = other.Cmd.Buffer.Strategy
	}
}
