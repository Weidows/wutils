package runner

import (
	"os"
	"path/filepath"
)

func ConfigTemplate() ([]byte, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(filepath.Dir(execPath), "config", "wutils.yml")
	return os.ReadFile(configPath)
}

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

type BufferConfig struct {
	Enable        bool   `default:"false"`
	MemoryLimit   int64  `default:"67108864"`
	FlushInterval int    `default:"10"`
	Strategy      string `default:"balanced"`
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
		Buffer BufferConfig
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
			Buffer BufferConfig
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
			Buffer: BufferConfig{Enable: false, MemoryLimit: 67108864, FlushInterval: 10, Strategy: "balanced"},
		},
	}
}
