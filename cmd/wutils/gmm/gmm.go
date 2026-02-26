package gmm

import (
	"fmt"
	"net/http"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"
)

type Mirror struct {
	Name string
	URL  string
}

var (
	ProxyMirrors = []Mirror{
		{"goproxy.cn", "https://goproxy.cn"},
		{"aliyun", "https://mirrors.aliyun.com/goproxy"},
		{"baidu", "https://goproxy.bj.bcebos.com"},
		{"huawei", "https://mirrors.huaweicloud.com/repository/goproxy"},
		{"tencent", "https://mirrors.cloud.tencent.com/go"},
		{"proxy-io", "https://goproxy.io"},
		{"tuna", "https://goproxy.cn"},
		{"default", "https://proxy.golang.org"},
	}

	SumdbMirrors = []Mirror{
		{"google", "https://sum.golang.google.cn"},
		{"sumdb-io", "https://gosum.io"},
		{"default", "https://sum.golang.org"},
	}
)

func TestSpeed() {
	fmt.Println("=== GOPROXY ===")
	testGroup(ProxyMirrors)
	fmt.Println("\n=== GOSUMDB ===")
	testGroup(SumdbMirrors)
}

func testGroup(mirrors []Mirror) {
	var wg sync.WaitGroup
	results := make(chan struct {
		name string
		ms   int64
	}, len(mirrors))

	for _, m := range mirrors {
		wg.Add(1)
		go func(m Mirror) {
			defer wg.Done()
			ms := measureLatency(m.URL)
			results <- struct {
				name string
				ms   int64
			}{m.Name, ms}
		}(m)
	}

	wg.Wait()
	close(results)

	sorted := make([]struct {
		name string
		ms   int64
	}, 0, len(mirrors))
	for r := range results {
		sorted = append(sorted, r)
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].ms == 0 {
			return false
		}
		if sorted[j].ms == 0 {
			return true
		}
		return sorted[i].ms < sorted[j].ms
	})

	for _, r := range sorted {
		if r.ms == 0 {
			fmt.Printf("timeout\t%s\n", r.name)
		} else {
			fmt.Printf("%dms\t%s\n", r.ms, r.name)
		}
	}
}

func measureLatency(url string) int64 {
	start := time.Now()
	client := &http.Client{Timeout: 3 * time.Second}
	_, err := client.Get(url + "/@v/list")
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		return 0
	}
	return elapsed
}

func SetProxy(name string) error {
	url := findURL(ProxyMirrors, name)
	if url == "" {
		return fmt.Errorf("unknown proxy: %s", name)
	}
	err := exec.Command("go", "env", "-w", "GOPROXY="+url+",direct").Run()
	if err == nil {
		fmt.Printf("GOPROXY=%s\n", url)
	}
	return err
}

func SetSumdb(name string) error {
	url := findURL(SumdbMirrors, name)
	if url == "" {
		return fmt.Errorf("unknown sumdb: %s", name)
	}
	cleanURL := strings.TrimPrefix(url, "https://")
	cleanURL = strings.TrimPrefix(cleanURL, "http://")
	err := exec.Command("go", "env", "-w", "GOSUMDB="+cleanURL).Run()
	if err == nil {
		fmt.Printf("GOSUMDB=%s\n", url)
	}
	return err
}

func findURL(mirrors []Mirror, name string) string {
	name = strings.ToLower(name)
	for _, m := range mirrors {
		if strings.ToLower(m.Name) == name {
			return m.URL
		}
	}
	return ""
}

func List() {
	fmt.Println("=== GOPROXY ===")
	for _, m := range ProxyMirrors {
		fmt.Printf("  %s\t%s\n", m.Name, m.URL)
	}
	fmt.Println("\n=== GOSUMDB ===")
	for _, m := range SumdbMirrors {
		fmt.Printf("  %s\t%s\n", m.Name, m.URL)
	}
}

func Current() {
	out, _ := exec.Command("go", "env", "GOPROXY").Output()
	fmt.Printf("GOPROXY=%s", string(out))
	out, _ = exec.Command("go", "env", "GOSUMDB").Output()
	fmt.Printf("GOSUMDB=%s", string(out))
}
