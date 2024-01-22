package net

import (
	time2 "github.com/Weidows/wutils/utils/time"
	"github.com/prometheus-community/pro-bing"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

/*
Ping returns Milliseconds (returns 0 if timeout by 3s)

	 1 Second(s)
		= 10e3 Milliseconds(ms)
		= 10e6 Microseconds(μs)
		= 10e9 Nanoseconds(ns)

`https://baidu.com` or `baidu.com` is OK.
*/
func Ping(host string) int64 {
	host = strings.TrimSpace(host)
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")

	ms := time2.WithTimeOut(time.Second*3, func() int64 {
		if runtime.GOOS == "windows" {
			return pingByLib(host)
		} else {
			return pingByHttp(host)
		}
	})
	return ms
}

func pingByLib(host string) int64 {
	p, err := probing.NewPinger(host)
	if err != nil {
		logger.Error(err)
	}

	p.Count = 3
	if runtime.GOOS == "windows" {
		//https://github.com/go-ping/ping#windows
		p.SetPrivileged(true)
	}

	if err = p.Run(); err != nil {
		logger.Error(err)
	}

	stats := p.Statistics()
	ms := stats.AvgRtt.Milliseconds()
	//logger.Println(stats.AvgRtt)
	return ms
}

func pingByHttp(host string) int64 {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	var resp *http.Response
	t := time2.TimeCosts(func() {
		resp, _ = client.Get("http://" + host)
	})
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return t.Milliseconds()
}

func NetReachable(host string) bool {
	cmd := exec.Command("ping", host, "-c", "4", "-W", "5")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}
