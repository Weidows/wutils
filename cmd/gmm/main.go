package main

import (
	"fmt"
	"github.com/Weidows/Golang/utils/collection"
	time2 "github.com/Weidows/Golang/utils/time"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// https://github.com/eryajf/Thanks-Mirror#go
var (
	PROXYS = map[string]map[string]string{
		"proxy": {
			"default":  "https://proxy.golang.org",
			"aliyun":   "https://mirrors.aliyun.com/goproxy",
			"proxy-cn": "https://goproxy.cn",
			"proxy-io": "https://proxy.golang.com.cn",
			"baidu":    "https://goproxy.bj.bcebos.com",
			"tencent":  "https://mirrors.cloud.tencent.com/go",
			"huawei":   "https://repo.huaweicloud.com/repository/goproxy",
		},
		"sumdb": {
			"default":  "https://sum.golang.org",
			"google":   "https://sum.golang.google.cn",
			"sumdb-io": "https://gosum.io",
		},
	}
	timeCost = map[string]map[int64]string{
		"proxy": make(map[int64]string),
		"sumdb": make(map[int64]string),
	}
	wg sync.WaitGroup
)

func main() {
	Commands := []*cli.Command{
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Speed test",
			Action: func(cCtx *cli.Context) (err error) {
				for k, v := range PROXYS {
					fmt.Println(k)
					for key, value := range v {
						wg.Add(1)
						go ping(k, value, key)
					}
					wg.Wait()
					sortedTimes := collection.SortKeys[int64](timeCost[k])
					for _, v1 := range sortedTimes {
						fmt.Printf("\t%dms\t%s\n", v1, timeCost[k][v1])
					}
				}
				return err
			},
		},
		{
			Name:    "proxy",
			Aliases: []string{"p"},
			Usage:   "change proxy",
			Action: func(cCtx *cli.Context) (err error) {
				input := strings.ToLower(cCtx.Args().First())
				s := PROXYS["proxy"][input]
				if s != "" {
					err = exec.Command("go", "env", "-w", "GOPROXY="+s+",direct").Run()
					fmt.Println("Proxy use "+input, s)
				} else {
					fmt.Println("Input name error: " + input)
				}
				return err
			},
		},
		{
			Name:    "sumdb",
			Aliases: []string{"s"},
			Usage:   "change sumdb",
			Action: func(cCtx *cli.Context) (err error) {
				input := strings.ToLower(cCtx.Args().First())
				// 不能带前面的 https, 会报错
				s := PROXYS["sumdb"][input][8:]
				if s != "" {
					err = exec.Command("go", "env", "-w", "GOSUMDB="+s).Run()
					fmt.Println("Sumdb use "+input, s)
				} else {
					fmt.Println("Input name error: " + input)
				}
				return err
			},
		},
	}

	app := &cli.App{
		Name: "Gmm",
		Usage: `
			You can test speed using 'gmm test',
			To change proxy command like 'gmm proxy huawei', same 'gmm sumdb google'
			Make sure 'GOPROXY/GOSUMDB' Not in system variable, otherwise gmm can's change mirrors.
		`,
		Commands: Commands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func ping(k string, value string, key string) {
	defer wg.Done()
	start := time.Now().UnixMilli()
	time2.WithTimeOut(time.Second*3, func() {
		_, _ = http.Get(value)
	})
	timeCost[k][time.Now().UnixMilli()-start] = key
}
