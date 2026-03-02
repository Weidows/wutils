package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Weidows/wutils/cmd/wutils/buffer"
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
				configPath = filepath.Join(home, ".config", "wutils", "app.yml")
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					templateData, err := runner.ConfigTemplate()
					if err != nil {
						return err
					}
					if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
						return err
					}
					err = os.WriteFile(configPath, templateData, 0644)
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
					{
						Name:  "update",
						Usage: "merge template config with user config, adding new fields while preserving user values",
						Action: func(cCtx *cli.Context) error {
							// Get template from executable directory
							templateData, err := runner.ConfigTemplate()
							if err != nil {
								return fmt.Errorf("failed to read template config: %v", err)
							}

							// Parse template as yaml.Node (preserves comments)
							var templateNode yaml.Node
							if err := yaml.Unmarshal(templateData, &templateNode); err != nil {
								return fmt.Errorf("failed to parse template: %v", err)
							}

							// Read user config
							userConfigData, err := os.ReadFile(configPath)
							if err != nil {
								return fmt.Errorf("failed to read user config: %v", err)
							}

							// Parse user config as yaml.Node (preserves comments)
							var userNode yaml.Node
							if err := yaml.Unmarshal(userConfigData, &userNode); err != nil {
								return fmt.Errorf("failed to parse user config: %v", err)
							}

							// Merge template with user config
							mergedNode := mergeYamlNodes(&templateNode, &userNode)

							// Encode merged result back to YAML
							mergedData, err := yaml.Marshal(mergedNode)
							if err != nil {
								return fmt.Errorf("failed to encode merged config: %v", err)
							}

							// Write merged config back to user config file
							if err := os.WriteFile(configPath, mergedData, 0644); err != nil {
								return fmt.Errorf("failed to write merged config: %v", err)
							}

							fmt.Printf("Config updated successfully: %s\n", configPath)
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
					if kr.Config.Cmd.Dsg.Parallel {
						go kr.Dsg()
					}
					if kr.Config.Cmd.Ol.Parallel {
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

			{
				Name:  "buffer",
				Usage: "Buffer filesystem operations",
				Subcommands: []*cli.Command{
					{
						Name:  "mount",
						Usage: "Mount a buffered drive",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "source",
								Aliases:  []string{"s"},
								Usage:    "source path to buffer (required)",
								Required: true,
							},
							&cli.Int64Flag{
								Name:  "memory-limit",
								Usage: "memory limit in bytes (default: 67108864)",
								Value: 67108864,
							},
							&cli.Int64Flag{
								Name:  "flush-interval",
								Usage: "flush interval in seconds (default: 10)",
								Value: 10,
							},
							&cli.StringFlag{
								Name:  "strategy",
								Usage: "buffer strategy: monitoring, defrag, download, migration, balanced (default: balanced)",
								Value: "balanced",
							},
							&cli.BoolFlag{
								Name:  "enable-read-cache",
								Usage: "enable read cache",
								Value: false,
							},
							&cli.BoolFlag{
								Name:  "enable-write-buffer",
								Usage: "enable write buffer",
								Value: true,
							},
						},
						Action: func(cCtx *cli.Context) error {
							drive := cCtx.Args().Get(0)
							if drive == "" {
								return fmt.Errorf("please provide drive letter (e.g., X:)")
							}

							// Get config from runner's Config.Cmd.Buffer
							cfg := kr.Config.Cmd.Buffer

							config := &buffer.BufferConfig{
								SourcePath:        cCtx.String("source"),
								MemoryLimit:       cCtx.Int64("memory-limit"),
								FlushInterval:     cCtx.Int64("flush-interval"),
								Strategy:          cCtx.String("strategy"),
								EnableReadCache:   cCtx.Bool("enable-read-cache"),
								EnableWriteBuffer: cCtx.Bool("enable-write-buffer"),
							}

							// Override with config file settings if not explicitly provided
							if config.MemoryLimit == 67108864 && cfg.MemoryLimit > 0 {
								config.MemoryLimit = cfg.MemoryLimit
							}
							if config.FlushInterval == 10 && cfg.FlushInterval > 0 {
								config.FlushInterval = int64(cfg.FlushInterval)
							}
							if config.Strategy == "balanced" && cfg.Strategy != "" {
								config.Strategy = cfg.Strategy
							}

							fmt.Printf("Mounting buffer drive %s from %s...\n", drive, config.SourcePath)
							if err := buffer.Mount(drive, config); err != nil {
								return fmt.Errorf("failed to mount buffer: %v", err)
							}
							fmt.Printf("Buffer drive %s mounted successfully\n", drive)
							return nil
						},
					},
					{
						Name:  "unmount",
						Usage: "Unmount the buffered drive",
						Action: func(cCtx *cli.Context) error {
							if err := buffer.Unmount(); err != nil {
								return fmt.Errorf("failed to unmount buffer: %v", err)
							}
							fmt.Println("Buffer unmounted successfully")
							return nil
						},
					},
					{
						Name:  "status",
						Usage: "Show buffer status",
						Action: func(cCtx *cli.Context) error {
							cfg := kr.Config.Cmd.Buffer
							fmt.Println("Buffer Configuration:")
							fmt.Printf("  Enable: %v\n", cfg.Enable)
							fmt.Printf("  MemoryLimit: %d bytes\n", cfg.MemoryLimit)
							fmt.Printf("  FlushInterval: %d seconds\n", cfg.FlushInterval)
							fmt.Printf("  Strategy: %s\n", cfg.Strategy)
							return nil
						},
					},
				},
			},
		},
	}
)

func mergeYamlNodes(template, user *yaml.Node) *yaml.Node {
	if template == nil {
		return user
	}
	if user == nil {
		return template
	}

	result := *template

	if template.Kind == yaml.MappingNode && user.Kind == yaml.MappingNode {
		templateMap := nodeToMap(template)
		userMap := nodeToMap(user)

		mergedContent := make([]*yaml.Node, 0)

		for i := 0; i < len(template.Content); i += 2 {
			key := template.Content[i]
			templateValue := template.Content[i+1]

			keyStr := key.Value
			if userValue, exists := userMap[keyStr]; exists {
				mergedContent = append(mergedContent, key)
				mergedContent = append(mergedContent, mergeYamlNodes(templateValue, userValue))
			} else {
				mergedContent = append(mergedContent, key)
				mergedContent = append(mergedContent, templateValue)
			}
		}

		for i := 0; i < len(user.Content); i += 2 {
			key := user.Content[i]
			keyStr := key.Value
			if _, exists := templateMap[keyStr]; !exists {
				mergedContent = append(mergedContent, key)
				mergedContent = append(mergedContent, user.Content[i+1])
			}
		}

		result.Content = mergedContent
	} else if user.Kind != 0 {
		return user
	}

	return &result
}

func nodeToMap(node *yaml.Node) map[string]*yaml.Node {
	result := make(map[string]*yaml.Node)
	if node.Kind != yaml.MappingNode {
		return result
	}

	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]
		if key.Kind == yaml.ScalarNode && key.Value != "" {
			result[key.Value] = value
		}
	}
	return result
}

func main() {
	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err.Error())
	}
}
