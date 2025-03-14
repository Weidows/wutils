package main

import (
	"fmt"
	"os"

	"github.com/Weidows/wutils/cmd/wutils/diff"
	"github.com/Weidows/wutils/cmd/wutils/keep_runner"
	"github.com/Weidows/wutils/cmd/wutils/media"
	"github.com/Weidows/wutils/cmd/wutils/zip"
	"github.com/Weidows/wutils/utils/log"
	"github.com/urfave/cli/v2"
)

var (
	logger = log.GetLogger()
	kr     = keep_runner.NewKeepRunner(logger)

	app = &cli.App{
		Name: "wutils",
		Authors: []*cli.Author{{
			Name:  "Weidows",
			Email: "ceo@weidows.tech",
		}},
		EnableBashCompletion: true,
		Usage: "Documents(使用指南) at here:\n" +
			"https://blog.weidows.tech/post/lang/golang/wutils",
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
					"文件对比工具, 但不是 Git-diff 那种\n" +
					"是用来求 '行-差集' 的工具\n" +
					"输入为两个特定名称的文件: './inputA.txt', './inputB.txt'",
				Action: func(cCtx *cli.Context) (err error) {
					missInA, missInB := diff.CheckLinesDiff("./inputA.txt", "./inputB.txt")
					// 输出结果
					fmt.Println("================== Missing in A ==================")
					for _, file := range missInA {
						fmt.Println(file)
					}

					fmt.Println("\n================== Missing in B:==================")
					for _, file := range missInB {
						fmt.Println(file)
					}

					return err
				},
			},

			{
				Name:    "parallel",
				Aliases: []string{"pl"},
				Usage:   "并行+后台执行任务 (配置取自wutils.yml)",
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
					"把指定窗口设置opacity, 使其透明化 (same as BLend)",
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
			{
				Name:  "zip",
				Usage: "some actions to operate zip/7z files",
				Subcommands: []*cli.Command{
					{
						Name:  "crack",
						Usage: "zip crack <path> - Crack the password of zip/7z files",
						Action: func(cCtx *cli.Context) error {
							if cCtx.Args().Len() < 1 {
								return fmt.Errorf("please provide the path to the archive file")
							}
							archivePath := cCtx.Args().Get(0)
							password := zip.CrackPassword(archivePath)

							if password == "" {
								return fmt.Errorf("no password found")
							} else {
								fmt.Printf("Password found: %s\n", password)
							}
							return nil
						},
					},
				},
			},

			{
				Name:  "media",
				Usage: "some actions to operate image or video files",
				Subcommands: []*cli.Command{
					{
						Name:  "group",
						Usage: "根据地理位置和时间对媒体文件进行分组",
						Action: func(c *cli.Context) error {
							inputDir := c.Args().Get(0)
							if inputDir == "" {
								return fmt.Errorf("请提供输入路径")
							}
							media.ClusterAndCopy(inputDir)

							return nil
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
