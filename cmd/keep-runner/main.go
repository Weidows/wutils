package main

import (
	"fmt"
	"github.com/Weidows/wutils/utils/collection"
	"github.com/Weidows/wutils/utils/files"
	"github.com/Weidows/wutils/utils/grammar"
	"github.com/Weidows/wutils/utils/log"
	os2 "github.com/Weidows/wutils/utils/os"
	"github.com/cheggaaa/pb/v3"
	"github.com/jinzhu/configor"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
	"time"
)

const ConfigPath = "keep-runner.yml"

var (
	logger = log.GetLogger()
	config = struct {
		Debug    bool `default:"false"`
		Parallel struct {
			Dsg bool
			Ol  bool
		}

		Dsg struct {
			Disk  []string `required:"true"`
			Delay int      `default:"30"`
		} `yaml:"dsg" required:"true"`

		Ol struct {
			Interval int `default:"1000"`
			Patterns []struct {
				Title   string
				Opacity byte
			}
		}
	}{}

	app = &cli.App{
		Name: "keep-runner",
		Authors: []*cli.Author{{
			Name:  "Weidows",
			Email: "utsuko27@gmail.com",
		}},
		EnableBashCompletion: true,
		Usage: "几个旨在后台运行的程序, config 使用: ./keep-runner.yml\n" +
			"Default config: https://github.com/Weidows/wutils/tree/master/config/cmd/keep-runner.yml",
		Commands: []*cli.Command{
			{
				Name:    "parallel",
				Aliases: []string{"pl"},
				Hidden:  true,
				Usage:   "并行+后台执行任务(取自config)",
				Action: func(cCtx *cli.Context) (err error) {
					parallel()
					return err
				},
			},
			{
				Name:      "dsg",
				Aliases:   []string{""},
				Hidden:    false,
				UsageText: "",
				Usage: "Disk sleep guard\n" +
					"防止硬盘睡眠 (每隔一段自定义的时间, 往指定盘里写一个时间戳)\n" +
					"外接 HDD 频繁启停甚是头疼, 后台让它怠速跑着, 免得起起停停增加损坏率",
				Action: func(cCtx *cli.Context) (err error) {
					dsg()
					return err
				},
			},
			{
				Name:    "ol",
				Aliases: []string{""},
				Hidden:  false,
				Usage: "Opacity Listener\n" +
					"后台持续运行, 并每隔指定时间扫一次运行的窗口\n" +
					"把指定窗口设置opacity, 使其透明化(比BLend好使~)",
				Action: func(cCtx *cli.Context) (err error) {
					ol()
					return err
				},
			},
			{
				Name:    "config",
				Aliases: []string{""},
				Hidden:  false,
				Usage:   "print config file",
				Action: func(cCtx *cli.Context) (err error) {
					logger.Println(fmt.Sprintf("%+v", config))
					return err
				},
			},
		},
	}
)

func parallel() {
	if config.Parallel.Dsg {
		go dsg()
	}
	if config.Parallel.Ol {
		ol()
	}

}

func dsg() {
	writeString := func(disk string) {
		f := strings.Join([]string{disk, ".dsg"}, "/")

		file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Println("disk format error, please input like 'E:'", err)
		}
		_, _ = file.WriteString("dsg running at " + time.Now().String() + "\n")
		_ = file.Close()
	}
	count := config.Dsg.Delay
	bar := pb.Simple.Start(count)

	for true {
		for i := 0; i < count; i++ {
			bar.Increment()
			time.Sleep(time.Second)
		}
		collection.ForEach(config.Dsg.Disk, func(i int, v string) {
			go writeString(v)
		})
		bar.SetCurrent(0)
	}
}

func ol() {
	count := config.Ol.Interval / 1000
	bar := pb.Simple.Start(count)

	for true {
		for i := 0; i < count; i++ {
			bar.Increment()
			time.Sleep(time.Second)
		}

		collection.ForEach(os2.GetEnumWindowsInfo(&os2.EnumWindowsFilter{
			IgnoreNoTitled:  true,
			IgnoreInvisible: true,
		}), func(i int, window *os2.EnumWindowsResult) {
			collection.ForEach(config.Ol.Patterns, func(ii int, pattern struct {
				Title   string
				Opacity byte
			}) {
				if grammar.Match(pattern.Title, window.Title) {
					isSuccess := os2.SetWindowOpacity(window.Handle, pattern.Opacity)
					if config.Debug {
						logger.Println(isSuccess, window, pattern)
					}
				}
			})
		})

		bar.SetCurrent(0)
	}
}

func main() {
	if !files.IsExist(ConfigPath) {
		return
	}
	_ = configor.Load(&config, ConfigPath)

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err.Error())
	}
}
