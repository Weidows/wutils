package keep_runner

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Weidows/wutils/utils/collection"
	"github.com/Weidows/wutils/utils/grammar"
	os2 "github.com/Weidows/wutils/utils/os"
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

const ConfigPath = "./Config/cmd/wutils.yml"

// Scope 定义了配置和日志记录器的结构体
type Scope struct {
	Logger *logrus.Logger
	Config struct {
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
}

// NewKeepRunner 初始化 Scope 实例
func NewKeepRunner(logger *logrus.Logger) *Scope {
	s := &Scope{Logger: logger}
	go s.init()
	time.Sleep(time.Millisecond * 50)
	return s
}

// init 加载配置文件并定期刷新
func (s *Scope) init() {
	for {
		_ = configor.Load(&s.Config, ConfigPath)
		if s.Config.Refresh.Delay < 1 {
			s.Config.Refresh.Delay = 1
		}
		time.Sleep(time.Second * time.Duration(s.Config.Refresh.Delay))
	}
}

// Dsg 执行 DSG 任务
func (s *Scope) Dsg() {
	s.Logger.Infoln(s.Config.Dsg)
	writeString := func(disk string) {
		f := strings.Join([]string{disk, ".dsg"}, "/")

		file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			s.Logger.Println("disk format error, please input like 'E:'", err)
		}
		_, _ = file.WriteString("dsg running at " + time.Now().String() + "\n")
		_ = file.Close()
	}

	for true {
		collection.ForEach(s.Config.Dsg.Disk, func(i int, v string) {
			go writeString(v)
		})
		time.Sleep(time.Second * time.Duration(s.Config.Dsg.Delay))
	}
}

// Ol 执行 OL 任务
func (s *Scope) Ol() {
	s.Logger.Infoln(s.Config.Ol)

	for true {
		windows := os2.GetEnumWindowsInfo(&os2.EnumWindowsFilter{
			IgnoreNoTitled:  true,
			IgnoreInvisible: true,
		})
		collection.ForEach(windows, func(i int, window *os2.EnumWindowsResult) {
			collection.ForEach(s.Config.Ol.Patterns, func(ii int, pattern struct {
				Title   string
				Opacity byte
			}) {
				if grammar.Match(pattern.Title, window.Title) && pattern.Opacity != window.Opacity {
					isSuccess := os2.SetWindowOpacity(window.Handle, pattern.Opacity)
					if s.Config.Debug {
						s.Logger.Println(isSuccess, window, pattern)
					}
				}
			})
		})
		time.Sleep(time.Second * time.Duration(s.Config.Ol.Delay))
	}
}

// OlList 列出所有符合条件的窗口信息
func (s *Scope) OlList() {
	collection.ForEach(os2.GetEnumWindowsInfo(&os2.EnumWindowsFilter{
		IgnoreNoTitled:  true,
		IgnoreInvisible: true,
	}), func(i int, v *os2.EnumWindowsResult) {
		// s.Logger.Println(fmt.Sprintf("%+v", v))
		fmt.Println(fmt.Sprintf("%+v", v))
	})
}
