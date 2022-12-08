package time

import (
	"testing"
	"time"
)

func TestWithTimeOut(t *testing.T) {
	var res string
	WithTimeOut(800*time.Millisecond, func() {
		time.Sleep(time.Millisecond)
		res += "1"

		time.Sleep(time.Second)
		res += "2"
	})
	println(res) // 1
}
