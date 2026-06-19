package runner

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Weidows/wutils/internal/config"
	"github.com/Weidows/wutils/utils/collection"
	"github.com/Weidows/wutils/utils/grammar"
	os2 "github.com/Weidows/wutils/utils/os"
	"github.com/sirupsen/logrus"
)

// Scope 定义了配置和日志记录器的结构体
type Scope struct {
	Logger     *logrus.Logger
	Config     *config.Config
	ConfigPath string
	watcher    *config.Watcher
}

// NewKeepRunner 初始化 Scope 实例
func NewKeepRunner(logger *logrus.Logger, configPath string) *Scope {
	s := &Scope{
		Logger:     logger,
		ConfigPath: configPath,
		Config:     loadConfig(configPath),
	}

	// Start hot-reload watcher
	interval := time.Duration(s.Config.Refresh) * time.Second
	if interval < time.Second {
		interval = time.Second
	}
	s.watcher = config.NewWatcher(configPath, interval)
	s.watcher.Start()

	return s
}

func loadConfig(path string) *config.Config {
	cfg, err := config.Load(path)
	if err != nil {
		cfg := config.DefaultConfig()
		return &cfg
	}
	return cfg
}

// Dsg 执行 DSG 任务
func (s *Scope) Dsg() {
	s.Logger.Infoln(s.Config.Cmd.Dsg)
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
		collection.ForEach(s.Config.Cmd.Dsg.Disk, func(i int, v string) {
			go writeString(v)
		})
		time.Sleep(time.Second * time.Duration(s.Config.Cmd.Dsg.Delay))
	}
}

// Ol 执行 OL 任务
func (s *Scope) Ol() {
	s.Logger.Infoln(s.Config.Cmd.Ol)

	for true {
		windows := os2.GetEnumWindowsInfo(&os2.EnumWindowsFilter{
			IgnoreNoTitled:  true,
			IgnoreInvisible: true,
		})
		collection.ForEach(windows, func(i int, window *os2.EnumWindowsResult) {
			collection.ForEach(s.Config.Cmd.Ol.Patterns, func(ii int, pattern config.OLPattern) {
				if grammar.Match(pattern.Title, window.Title) && pattern.Opacity != window.Opacity {
					isSuccess := os2.SetWindowOpacity(window.Handle, pattern.Opacity)
					if s.Config.App.Debug {
						s.Logger.Println(isSuccess, window, pattern)
					}
				}
			})
		})
		time.Sleep(time.Second * time.Duration(s.Config.Cmd.Ol.Delay))
	}
}

// OlList 列出所有符合条件的窗口信息
func (s *Scope) OlList() {
	collection.ForEach(os2.GetEnumWindowsInfo(&os2.EnumWindowsFilter{
		IgnoreNoTitled:  true,
		IgnoreInvisible: true,
	}), func(i int, v *os2.EnumWindowsResult) {
		fmt.Println(fmt.Sprintf("%+v", v))
	})
}
