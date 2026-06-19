package service

import (
	"fmt"
	"net/http"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/i18n"
)

type mirror struct {
	name string
	url  string
}

var (
	proxyMirrors = []mirror{
		{"goproxy.cn", "https://goproxy.cn"},
		{"aliyun", "https://mirrors.aliyun.com/goproxy"},
		{"baidu", "https://goproxy.bj.bcebos.com"},
		{"huawei", "https://mirrors.huaweicloud.com/repository/goproxy"},
		{"tencent", "https://mirrors.cloud.tencent.com/go"},
		{"proxy-io", "https://goproxy.io"},
		{"tuna", "https://goproxy.cn"},
		{"default", "https://proxy.golang.org"},
	}
	sumdbMirrors = []mirror{
		{"google", "https://sum.golang.google.cn"},
		{"sumdb-io", "https://gosum.io"},
		{"default", "https://sum.golang.org"},
	}
)

// GMMService manages Go module proxy mirrors (GOPROXY, GOSUMDB).
type GMMService struct{}

// NewGMMService creates a new GMMService.
func NewGMMService() *GMMService {
	return &GMMService{}
}

func (s *GMMService) Name() string               { return "gmm" }
func (s *GMMService) Description() string        { return i18n.G("gmm.description") }
func (s *GMMService) Status() app.ServiceStatus  { return app.StatusStopped }
func (s *GMMService) Start() error               { return nil }
func (s *GMMService) Stop() error                { return nil }

// TestSpeed benchmarks all proxy mirrors.
func (s *GMMService) TestSpeed() {
	fmt.Println("=== GOPROXY ===")
	testGroup(proxyMirrors)
	fmt.Println("\n=== GOSUMDB ===")
	testGroup(sumdbMirrors)
}

func testGroup(mirrors []mirror) {
	var wg sync.WaitGroup
	type result struct {
		name string
		ms   int64
	}
	ch := make(chan result, len(mirrors))

	for _, m := range mirrors {
		wg.Add(1)
		go func(m mirror) {
			defer wg.Done()
			ms := measureLatency(m.url)
			ch <- result{m.name, ms}
		}(m)
	}
	wg.Wait()
	close(ch)

	var sorted []result
	for r := range ch {
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
	if err != nil {
		return 0
	}
	return time.Since(start).Milliseconds()
}

// SetProxy configures GOPROXY to the named mirror.
func (s *GMMService) SetProxy(name string) error {
	url := findURL(proxyMirrors, name)
	if url == "" {
		return fmt.Errorf("unknown proxy: %s", name)
	}
	if err := exec.Command("go", "env", "-w", "GOPROXY="+url+",direct").Run(); err != nil {
		return err
	}
	fmt.Printf("GOPROXY=%s\n", url)
	return nil
}

// SetSumdb configures GOSUMDB to the named mirror.
func (s *GMMService) SetSumdb(name string) error {
	url := findURL(sumdbMirrors, name)
	if url == "" {
		return fmt.Errorf("unknown sumdb: %s", name)
	}
	cleanURL := strings.TrimPrefix(url, "https://")
	cleanURL = strings.TrimPrefix(cleanURL, "http://")
	if err := exec.Command("go", "env", "-w", "GOSUMDB="+cleanURL).Run(); err != nil {
		return err
	}
	fmt.Printf("GOSUMDB=%s\n", url)
	return nil
}

func findURL(mirrors []mirror, name string) string {
	name = strings.ToLower(name)
	for _, m := range mirrors {
		if strings.ToLower(m.name) == name {
			return m.url
		}
	}
	return ""
}

// List prints all available mirrors.
func (s *GMMService) List() {
	fmt.Println("=== GOPROXY ===")
	for _, m := range proxyMirrors {
		fmt.Printf("  %s\t%s\n", m.name, m.url)
	}
	fmt.Println("\n=== GOSUMDB ===")
	for _, m := range sumdbMirrors {
		fmt.Printf("  %s\t%s\n", m.name, m.url)
	}
}

// Current prints the current go environment proxy configuration.
func (s *GMMService) Current() {
	out, _ := exec.Command("go", "env", "GOPROXY").Output()
	fmt.Printf("GOPROXY=%s", string(out))
	out, _ = exec.Command("go", "env", "GOSUMDB").Output()
	fmt.Printf("GOSUMDB=%s", string(out))
}
