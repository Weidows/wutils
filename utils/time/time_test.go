package time

import (
	"testing"
	"time"
)

func TestWithTimeOut(t *testing.T) {
	logger.Println(WithTimeOut(800*time.Millisecond, func() int {
		time.Sleep(time.Millisecond)
		return 1
	})) // 1

	logger.Println(WithTimeOut(800*time.Millisecond, func() string {
		time.Sleep(time.Second)
		return "2"
	})) // nil
}
