package main

import (
	"fmt"
	"github.com/Weidows/wutils/utils/collection"
	"github.com/Weidows/wutils/utils/grammar"
	"github.com/Weidows/wutils/utils/log"
	os2 "github.com/Weidows/wutils/utils/os"
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
				Usage:   "并行+后台执行任务(取自config)",
				Action: func(cCtx *cli.Context) (err error) {
					if config.Parallel.Dsg {
						go dsg()
					}
					if config.Parallel.Ol {
						ol()
					}

					return err
				},
			},
			{
				Name:      "dsg",
				Aliases:   []string{""},
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
				Usage: "Opacity Listener\n" +
					"后台持续运行, 并每隔指定时间扫一次运行的窗口\n" +
					"把指定窗口设置opacity, 使其透明化(比BLend好使~)",
				Action: func(cCtx *cli.Context) (err error) {
					ol()
					return err
				},
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{""},
						Usage:   "list all visible windows",
						Action: func(cCtx *cli.Context) (err error) {
							olList()
							return err
						},
					},
				},
			},
			{
				Name:    "config",
				Aliases: []string{""},
				Usage:   "print config file",
				Action: func(cCtx *cli.Context) (err error) {
					logger.Println(fmt.Sprintf("%+v", config))
					return err
				},
			},
		},
	}
)

func dsg() {
	logger.Infoln(config.Dsg)
	writeString := func(disk string) {
		f := strings.Join([]string{disk, ".dsg"}, "/")

		file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Println("disk format error, please input like 'E:'", err)
		}
		_, _ = file.WriteString("dsg running at " + time.Now().String() + "\n")
		_ = file.Close()
	}

	for true {
		collection.ForEach(config.Dsg.Disk, func(i int, v string) {
			go writeString(v)
		})
		time.Sleep(time.Second * time.Duration(config.Dsg.Delay))
	}
}

func ol() {
	logger.Infoln(config.Ol)

	for true {
		windows := os2.GetEnumWindowsInfo(&os2.EnumWindowsFilter{
			IgnoreNoTitled:  true,
			IgnoreInvisible: true,
		})
		collection.ForEach(windows, func(i int, window *os2.EnumWindowsResult) {
			collection.ForEach(config.Ol.Patterns, func(ii int, pattern struct {
				Title   string
				Opacity byte
			}) {
				if grammar.Match(pattern.Title, window.Title) && pattern.Opacity != window.Opacity {
					//logger.Println(window, pattern)
					isSuccess := os2.SetWindowOpacity(window.Handle, pattern.Opacity)
					if config.Debug {
						logger.Println(isSuccess, window, pattern)
					}
				}
			})
		})
		time.Sleep(time.Second * time.Duration(config.Ol.Delay))
	}
}

func olList() {
	collection.ForEach(os2.GetEnumWindowsInfo(&os2.EnumWindowsFilter{
		IgnoreNoTitled:  true,
		IgnoreInvisible: true,
	}), func(i int, v *os2.EnumWindowsResult) {
		logger.Println(fmt.Sprintf("%+v", v))
	})
}

func refreshConfig() {
	for {
		_ = configor.Load(&config, ConfigPath)
		time.Sleep(time.Second * time.Duration(config.Refresh.Delay))
	}
}

func main() {
	go refreshConfig()
	time.Sleep(500)

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err.Error())
	}
}
