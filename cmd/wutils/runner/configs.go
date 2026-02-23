package runner

// Config 定义配置结构体
type Config struct {
	Debug   bool `default:"false"`
	Refresh struct {
		Delay int `default:"10"`
	}
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
