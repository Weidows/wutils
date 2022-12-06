package net

import (
	"github.com/go-ping/ping"
	"runtime"
)

/*
Ping returns Milliseconds

	 1 Second(s)
		= 10e3 Milliseconds(ms)
		= 10e6 Microseconds(μs)
		= 10e9 Nanoseconds(ns)
*/
func Ping(host string) int64 {
	p, err := ping.NewPinger(host)
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
	return stats.AvgRtt.Milliseconds()
}
