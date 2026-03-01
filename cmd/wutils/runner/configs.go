package runner

type AppConfig struct {
	Name    string `default:"wutils"`
	Version string `default:"1.0.0"`
	Debug   bool   `default:"false"`
}

type LoggingConfig struct {
	Level  string `default:"info"`
	Format string `default:"json"`
}

type RunnerConfig struct {
	MaxWorkers int    `default:"4"`
	Timeout    string `default:"30s"`
	Refresh    int    `default:"10"`
}

type Config struct {
	App     AppConfig
	Logging LoggingConfig
	Refresh int `default:"10"`
	Cmd     struct {
		Dsg struct {
			Parallel bool
			Disk     []string `required:"true"`
			Delay    int      `default:"30"`
		}
		Ol struct {
			Parallel bool
			Delay    int `default:"2"`
			Patterns []struct {
				Title   string
				Opacity byte
			}
		}
	}
}

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
		Cmd: struct {
			Dsg struct {
				Parallel bool
				Disk     []string `required:"true"`
				Delay    int      `default:"30"`
			}
			Ol struct {
				Parallel bool
				Delay    int `default:"2"`
				Patterns []struct {
					Title   string
					Opacity byte
				}
			}
		}{
			Dsg: struct {
				Parallel bool
				Disk     []string `required:"true"`
				Delay    int      `default:"30"`
			}{Parallel: true, Disk: []string{"E:"}, Delay: 30},
			Ol: struct {
				Parallel bool
				Delay    int `default:"2"`
				Patterns []struct {
					Title   string
					Opacity byte
				}
			}{Parallel: true, Delay: 2, Patterns: []struct {
				Title   string
				Opacity byte
			}{
				{Title: "(XY|xy)plorer", Opacity: 200},
				{Title: "设置$", Opacity: 220},
			}},
		},
	}
}
