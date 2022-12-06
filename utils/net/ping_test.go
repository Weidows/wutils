package net

import "testing"

func TestPing(t *testing.T) {
	p := Ping("baidu.com")
	logger.Println(p)
}
