// Package cli builds the urfave/cli application, wiring commands to the service layer.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Weidows/wutils/cmd/wutils/buffer"
	"github.com/Weidows/wutils/internal/config"
	"github.com/Weidows/wutils/internal/service"
	"github.com/Weidows/wutils/utils/log"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

// staticProvider wraps a fixed Config as a ConfigProvider (no hot-reload).
type staticProvider struct {
	cfg *config.Config
}

func (p *staticProvider) Get() *config.Config { return p.cfg }

// NewApp builds and returns the wutils CLI application.
func NewApp() *cli.App {
	logger := log.GetLogger()

	// Shared state
	var (
		cfgPath  string
		cfgProv  config.ConfigProvider
		watcher  *config.Watcher
	)

	// Services (lazily initialized)
	var (
		dsgSvc    *service.DSGService
		olSvc     *service.OLService
		diffSvc   = service.NewDiffService()
		zipSvc    = service.NewZipCrackService()
		mediaSvc  = service.NewMediaService()
		extractSvc = service.NewExtractService()
		gmmSvc    = service.NewGMMService()
	)

	return &cli.App{
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
			cfgPath = cCtx.String("config")
			if cfgPath == "" {
				home, _ := os.UserHomeDir()
				cfgPath = filepath.Join(home, ".config", "wutils", "app.yml")
				if err := config.EnsureUserConfig(cfgPath); err != nil {
					return fmt.Errorf("failed to ensure config: %w", err)
				}
			}

			// Load config
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Start config watcher for hot-reload
			interval := cfg.Refresh
			if interval < 1 {
				interval = 1
			}
			watcher = config.NewWatcher(cfgPath, time.Duration(cfg.Refresh)*time.Second)
			watcher.Start()

			// Use watcher as the live config provider for services
			cfgProv = watcher

			dsgSvc = service.NewDSGService(cfgProv, logger)
			olSvc = service.NewOLService(cfgProv, logger)

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
						Action: func(cCtx *cli.Context) error {
							fmt.Printf("Config file: %s\n\n", cfgPath)
							content, err := os.ReadFile(cfgPath)
							if err != nil {
								return fmt.Errorf("failed to read config: %v", err)
							}
							fmt.Println(string(content))
							return nil
						},
					},
					{
						Name:  "update",
						Usage: "merge template config with user config, preserving user values",
						Action: func(cCtx *cli.Context) error {
							templateData, err := config.Template()
							if err != nil {
								return fmt.Errorf("failed to read template: %v", err)
							}
							var templateNode yaml.Node
							if err := yaml.Unmarshal(templateData, &templateNode); err != nil {
								return fmt.Errorf("failed to parse template: %v", err)
							}
							userData, err := os.ReadFile(cfgPath)
							if err != nil {
								return fmt.Errorf("failed to read user config: %v", err)
							}
							var userNode yaml.Node
							if err := yaml.Unmarshal(userData, &userNode); err != nil {
								return fmt.Errorf("failed to parse user config: %v", err)
							}
							merged := config.MergeYAMLNodes(&templateNode, &userNode)
							mergedData, err := yaml.Marshal(merged)
							if err != nil {
								return fmt.Errorf("failed to encode merged config: %v", err)
							}
							if err := os.WriteFile(cfgPath, mergedData, 0644); err != nil {
								return fmt.Errorf("failed to write config: %v", err)
							}
							fmt.Printf("Config updated: %s\n", cfgPath)
							return nil
						},
					},
				},
			},

			{
				Name: "diff",
				Usage: "diff — 文件行差集对比\n" +
					"求两个文件的行差集（对称差），输入为 inputA.txt 和 inputB.txt",
				Action: func(cCtx *cli.Context) error {
					missInA, missInB := diffSvc.Diff("./inputA.txt", "./inputB.txt")
					fmt.Println("================== Missing in A ==================")
					for _, f := range missInA {
						fmt.Println(f)
					}
					fmt.Println("\n================== Missing in B:==================")
					for _, f := range missInB {
						fmt.Println(f)
					}
					return nil
				},
			},

			{
				Name:    "parallel",
				Aliases: []string{"pl"},
				Usage:   "并行+后台执行任务 (配置取自 wutils.yml)",
				Action: func(cCtx *cli.Context) error {
					_ = dsgSvc.Start()
					_ = olSvc.Start()
					return nil
				},
			},

			{
				Name: "dsg",
				Usage: "Disk Sleep Guard — 防止硬盘睡眠\n" +
					"在指定磁盘上定时写入时间戳，防止外接 HDD 频繁启停",
				Action: func(cCtx *cli.Context) error {
					return dsgSvc.Start()
				},
			},

			{
				Name: "ol",
				Usage: "Opacity Listener — 窗口透明度控制\n" +
					"后台持续扫描窗口，根据匹配规则自动设置透明度",
				Action: func(cCtx *cli.Context) error {
					return olSvc.Start()
				},
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "列出所有可见窗口",
						Action: func(cCtx *cli.Context) error {
							for _, w := range olSvc.ListWindows() {
								fmt.Printf("%+v\n", w)
							}
							return nil
						},
					},
				},
			},

			{
				Name:  "zip",
				Usage: "压缩文件操作",
				Subcommands: []*cli.Command{
					{
						Name:  "crack",
						Usage: "zip crack <path> — 破解压缩包密码",
						Action: func(cCtx *cli.Context) error {
							if cCtx.Args().Len() < 1 {
								return fmt.Errorf("请提供压缩包路径")
							}
							password := zipSvc.CrackPassword(cCtx.Args().Get(0))
							if password == "" {
								return fmt.Errorf("未找到密码")
							}
							fmt.Printf("Password found: %s\n", password)
							return nil
						},
					},
				},
			},

			{
				Name:  "media",
				Usage: "媒体文件操作",
				Subcommands: []*cli.Command{
					{
						Name:  "group",
						Usage: "根据地理位置和时间对媒体文件进行分组",
						Action: func(c *cli.Context) error {
							inputDir := c.Args().Get(0)
							if inputDir == "" {
								return fmt.Errorf("请提供输入路径")
							}
							mediaSvc.ClusterAndCopy(inputDir)
							return nil
						},
					},
				},
			},

			{
				Name:  "extract",
				Usage: "解散一级目录 — 将子目录内容提取到父目录",
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 2 {
						return fmt.Errorf("用法: wutils extract <mode> <path>\nmode: 0=自动检查, 1=覆盖, 2=跳过")
					}
					return extractSvc.Run(c.Args().Get(0), c.Args().Get(1))
				},
			},

			{
				Name:  "gmm",
				Usage: "Go Mirror Manager — 测试/切换模块代理",
				Subcommands: []*cli.Command{
					{
						Name:    "test",
						Aliases: []string{"t"},
						Usage:   "测试所有代理速度",
						Action:  func(c *cli.Context) error { gmmSvc.TestSpeed(); return nil },
					},
					{
						Name:    "proxy",
						Aliases: []string{"p"},
						Usage:   "设置 GOPROXY",
						Action: func(c *cli.Context) error {
							if c.Args().Len() < 1 {
								return fmt.Errorf("用法: wutils gmm proxy <name>")
							}
							return gmmSvc.SetProxy(c.Args().First())
						},
					},
					{
						Name:    "sumdb",
						Aliases: []string{"s"},
						Usage:   "设置 GOSUMDB",
						Action: func(c *cli.Context) error {
							if c.Args().Len() < 1 {
								return fmt.Errorf("用法: wutils gmm sumdb <name>")
							}
							return gmmSvc.SetSumdb(c.Args().First())
						},
					},
					{
						Name:  "list",
						Usage: "列出所有可用代理",
						Action: func(c *cli.Context) error { gmmSvc.List(); return nil },
					},
					{
						Name:  "current",
						Usage: "显示当前代理配置",
						Action: func(c *cli.Context) error { gmmSvc.Current(); return nil },
					},
				},
			},

			{
				Name:  "buffer",
				Usage: "IO 缓冲虚拟文件系统 (基于 Dokan)",
				Subcommands: []*cli.Command{
					{
						Name:  "mount",
						Usage: "挂载缓冲盘",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "source",
								Aliases:  []string{"s"},
								Usage:    "源路径 (必填)",
								Required: true,
							},
							&cli.Int64Flag{
								Name:  "memory-limit",
								Usage: "内存限制 (bytes, 默认 67108864)",
								Value: 67108864,
							},
							&cli.Int64Flag{
								Name:  "flush-interval",
								Usage: "刷新间隔 (秒, 默认 10)",
								Value: 10,
							},
							&cli.StringFlag{
								Name:  "strategy",
								Usage: "策略: monitoring/defrag/download/migration/balanced",
								Value: "balanced",
							},
							&cli.BoolFlag{
								Name:  "enable-read-cache",
								Usage: "启用读缓存",
							},
							&cli.BoolFlag{
								Name:  "enable-write-buffer",
								Usage: "启用写缓冲",
								Value: true,
							},
						},
						Action: func(cCtx *cli.Context) error {
							drive := cCtx.Args().Get(0)
							if drive == "" {
								return fmt.Errorf("请提供盘符 (如 X:)")
							}
							bufCfg := &buffer.BufferConfig{
								SourcePath:        cCtx.String("source"),
								MemoryLimit:       cCtx.Int64("memory-limit"),
								FlushInterval:     cCtx.Int64("flush-interval"),
								Strategy:          cCtx.String("strategy"),
								EnableReadCache:   cCtx.Bool("enable-read-cache"),
								EnableWriteBuffer: cCtx.Bool("enable-write-buffer"),
							}
							bufSvc := service.NewBufferService(drive, bufCfg)
							return bufSvc.Start()
						},
					},
					{
						Name:  "unmount",
						Usage: "卸载缓冲盘",
						Action: func(cCtx *cli.Context) error {
							bufSvc := service.NewBufferService("", nil)
							return bufSvc.Stop()
						},
					},
					{
						Name:  "status",
						Usage: "显示缓冲状态",
						Action: func(cCtx *cli.Context) error {
							cfg := watcher.Get()
							fmt.Println("Buffer Configuration:")
							fmt.Printf("  Enable: %v\n", cfg.Cmd.Buffer.Enable)
							fmt.Printf("  MemoryLimit: %d bytes\n", cfg.Cmd.Buffer.MemoryLimit)
							fmt.Printf("  FlushInterval: %d seconds\n", cfg.Cmd.Buffer.FlushInterval)
							fmt.Printf("  Strategy: %s\n", cfg.Cmd.Buffer.Strategy)
							return nil
						},
					},
				},
			},
		},
	}
}
