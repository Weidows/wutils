package net

import (
	"testing"
)

func TestPing(t *testing.T) {
	p := Ping("github.com")
	logger.Println(p)
}
