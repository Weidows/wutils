package time

import (
	"testing"
	"time"
)

func TestWithTimeOut(t *testing.T) {
	spendTime()
}

func spendTime() {
	logger.Println(WithTimeOut(800*time.Millisecond, func() int {
		time.Sleep(time.Millisecond * 500)
		return 1
	})) // 1

	logger.Println(WithTimeOut(800*time.Millisecond, func() int64 {
		time.Sleep(time.Second)
		return 2
	})) // 0

	logger.Println(WithTimeOut(800*time.Millisecond, func() string {
		time.Sleep(time.Second)
		return "3"
	})) // ""
}

func TestTimeCosts(t *testing.T) {
	logger.Println(TimeCosts(spendTime))
}
