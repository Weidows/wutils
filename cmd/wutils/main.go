package main

import (
	"fmt"
	"os"

	"github.com/Weidows/wutils/cmd/wutils/diff"
	"github.com/Weidows/wutils/cmd/wutils/keep_runner"
	"github.com/Weidows/wutils/utils/log"
	"github.com/urfave/cli/v2"
)

var (
	logger = log.GetLogger()
	kr     = keep_runner.NewKeepRunner(logger)

	app = &cli.App{
		Name: "keep-runner",
		Authors: []*cli.Author{{
			Name:  "Weidows",
			Email: "ceo@weidows.tech",
		}},
		EnableBashCompletion: true,
		Usage: "几个旨在后台运行的程序, config 使用: ./keep-runner.yml\n" +
			"Default config: https://github.com/Weidows/wutils/tree/master/config/cmd/keep-runner.yml",
		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: "print config file",
				Action: func(cCtx *cli.Context) (err error) {
					logger.Println(fmt.Sprintf("%+v", kr.Config))
					return err
				},
			},

			{
				Name: "diff",
				Usage: "diff - Differential set between two files\n" +
					"文件对比工具, 但不是 Git-diff 那种, 是用来求 '行-差集' 的工具",
				Action: func(cCtx *cli.Context) (err error) {
					missInA, missInB := diff.CheckLinesDiff("./inputA.txt", "./inputB.txt")
					// 输出结果
					fmt.Println("Missing in A")
					for _, file := range missInA {
						fmt.Println(file)
					}

					fmt.Println("\nMissing in B")
					for _, file := range missInB {
						fmt.Println(file)
					}

					return err
				},
			},

			{
				Name:    "parallel",
				Aliases: []string{"pl"},
				Usage:   "并行+后台执行任务(取自config)",
				Action: func(cCtx *cli.Context) (err error) {
					if kr.Config.Parallel.Dsg {
						go kr.Dsg()
					}
					if kr.Config.Parallel.Ol {
						kr.Ol()
					}

					return err
				},
			},
			{
				Name:      "dsg",
				UsageText: "",
				Usage: "Disk sleep guard\n" +
					"防止硬盘睡眠 (每隔一段自定义的时间, 往指定盘里写一个时间戳)\n" +
					"外接 HDD 频繁启停甚是头疼, 后台让它怠速跑着, 免得起起停停增加损坏率",
				Action: func(cCtx *cli.Context) (err error) {
					kr.Dsg()
					return err
				},
			},
			{
				Name: "ol",
				Usage: "Opacity Listener\n" +
					"后台持续运行, 并每隔指定时间扫一次运行的窗口\n" +
					"把指定窗口设置opacity, 使其透明化(比BLend好使~)",
				Action: func(cCtx *cli.Context) (err error) {
					kr.Ol()
					return err
				},
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "list all visible windows",
						Action: func(cCtx *cli.Context) (err error) {
							kr.OlList()
							return err
						},
					},
				},
			},
		},
	}
)

func main() {
	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err.Error())
	}
}
