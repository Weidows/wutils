package main

import (
	"testing"
	"time"
)

func Test_ol(t *testing.T) {
	go refreshConfig()
	time.Sleep(500)

	//ol()
	olList()
}
