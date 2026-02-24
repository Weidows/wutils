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
	App      AppConfig
	Logging  LoggingConfig
	Runner   RunnerConfig
	Parallel struct {
		Dsg bool
		Ol  bool
	}
	Dsg struct {
		Disk  []string `required:"true"`
		Delay int      `default:"30"`
	} `yaml:"dsg" required:"true"`
	Ol struct {
		Delay    int `default:"2"`
		Patterns []struct {
			Title   string
			Opacity byte
		}
	}
}
