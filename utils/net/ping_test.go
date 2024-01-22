package net

import (
	"testing"
)

func TestPing(t *testing.T) {
	logger.Println(Ping("baidu.com"))
	logger.Println(Ping("github.com"))
}
