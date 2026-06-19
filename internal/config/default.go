package config

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		App: AppConfig{
			Name:    "wutils",
			Version: "1.0.0",
			Debug:   false,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Refresh: 10,
		Cmd: CmdConfig{
			Dsg: DSGConfig{
				Parallel: true,
				Disk:     []string{"E:"},
				Delay:    30,
			},
			Ol: OLConfig{
				Parallel: true,
				Delay:    2,
				Patterns: []OLPattern{
					{Title: "(XY|xy)plorer", Opacity: 200},
					{Title: "设置$", Opacity: 220},
				},
			},
			Buffer: BufferConfig{
				Enable:        false,
				MemoryLimit:   67108864,
				FlushInterval: 10,
				Strategy:      "balanced",
			},
		},
	}
}
