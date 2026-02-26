package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Weidows/wutils/cmd/wutils/diff"
	"github.com/Weidows/wutils/cmd/wutils/extract"
	"github.com/Weidows/wutils/cmd/wutils/gmm"
	"github.com/Weidows/wutils/cmd/wutils/media"
	"github.com/Weidows/wutils/cmd/wutils/runner"
	"github.com/Weidows/wutils/cmd/wutils/zip"
	logutil "github.com/Weidows/wutils/utils/log"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var (
	logger     = logutil.GetLogger()
	kr         *runner.Scope
	configPath string

	app = &cli.App{
		Name: "wutils",
		Authors: []*cli.Author{{
			Name:  "Weidows",
			Email: "ceo@weidows.tech",
		}},
		EnableBashCompletion: true,
		Usage: "Documents(使用指南) at here:\n" +
			"https://blog.weidows.tech/post/lang/golang/wutils",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "path to config file",
			},
		},
		Before: func(cCtx *cli.Context) error {
			configPath = cCtx.String("config")
			if configPath == "" {
				home, _ := os.UserHomeDir()
				ymlPath := filepath.Join(home, ".wutils.yml")
				yamlPath := filepath.Join(home, ".wutils.yaml")
				if _, err := os.Stat(ymlPath); err == nil {
					configPath = ymlPath
				} else if _, err := os.Stat(yamlPath); err == nil {
					configPath = yamlPath
				} else {
					// Generate default config at ~/.wutils.yml
					configPath = ymlPath
					defaultConfig := runner.Config{
						App: runner.AppConfig{
							Name:    "wutils",
							Version: "1.0.0",
							Debug:   false,
						},
						Logging: runner.LoggingConfig{
							Level:  "info",
							Format: "json",
						},
						Runner: runner.RunnerConfig{
							MaxWorkers: 4,
							Timeout:    "30s",
						},
						Parallel: struct {
							Dsg bool
							Ol  bool
						}{Dsg: true, Ol: true},
						Dsg: struct {
							Disk  []string `required:"true"`
							Delay int      `default:"30"`
						}{Disk: []string{"E:"}, Delay: 30},
						Ol: struct {
							Delay    int `default:"2"`
							Patterns []struct {
								Title   string
								Opacity byte
							}
						}{Delay: 2, Patterns: []struct {
							Title   string
							Opacity byte
						}{
							{Title: "(XY|xy)plorer", Opacity: 200},
							{Title: "设置$", Opacity: 220},
						}},
					}
					data, err := yaml.Marshal(&defaultConfig)
					if err != nil {
						return err
					}
					if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
						return err
					}
					err = os.WriteFile(configPath, data, 0644)
					if err != nil {
						return err
					}
				}
			}
			kr = runner.NewKeepRunner(logger, configPath)
			return nil
		},

		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: "config operations",
				Subcommands: []*cli.Command{
					{
						Name:  "cat",
						Usage: "show config file location and content",
						Action: func(cCtx *cli.Context) (err error) {
							fmt.Printf("Config file: %s\n\n", configPath)
							content, err := os.ReadFile(configPath)
							if err != nil {
								return fmt.Errorf("failed to read config file: %v", err)
							}
							fmt.Println(string(content))
							return nil
						},
					},
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

			{
				Name:  "extract",
				Usage: "解散一级目录, 将子目录内容提取到父目录",
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 2 {
						return fmt.Errorf("用法: wutils extract <mode> <path>\nmode: 0=自动检查, 1=覆盖, 2=跳过")
					}
					mode := c.Args().Get(0)
					rootPath := c.Args().Get(1)
					return extract.Run(mode, rootPath)
				},
			},

			{
				Name:  "gmm",
				Usage: "Golang Mirror Manager - 测试/切换 Go 模块代理",
				Subcommands: []*cli.Command{
					{
						Name:    "test",
						Aliases: []string{"t"},
						Usage:   "测试所有代理速度",
						Action: func(c *cli.Context) error {
							gmm.TestSpeed()
							return nil
						},
					},
					{
						Name:    "proxy",
						Aliases: []string{"p"},
						Usage:   "设置 GOPROXY",
						Action: func(c *cli.Context) error {
							if c.Args().Len() < 1 {
								return fmt.Errorf("用法: wutils gmm proxy <name>\n可用: aliyun, baidu, huawei, tencent, goproxy.cn, proxy-io, default")
							}
							return gmm.SetProxy(c.Args().First())
						},
					},
					{
						Name:    "sumdb",
						Aliases: []string{"s"},
						Usage:   "设置 GOSUMDB",
						Action: func(c *cli.Context) error {
							if c.Args().Len() < 1 {
								return fmt.Errorf("用法: wutils gmm sumdb <name>\n可用: default, google, sumdb-io")
							}
							return gmm.SetSumdb(c.Args().First())
						},
					},
					{
						Name:  "list",
						Usage: "列出所有可用代理",
						Action: func(c *cli.Context) error {
							gmm.List()
							return nil
						},
					},
					{
						Name:  "current",
						Usage: "显示当前代理配置",
						Action: func(c *cli.Context) error {
							gmm.Current()
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
