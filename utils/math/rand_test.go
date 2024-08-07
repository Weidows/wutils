package math

import (
	"fmt"
	"testing"
	"time"
)

func Test_GetRandNum(t *testing.T) {
	for i := 0; i < 30; i++ {
		time.Sleep(time.Nanosecond)
		fmt.Println(GetRandNum(i))
	}
}

func Test_GetVerifyCode(t *testing.T) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Nanosecond)
		fmt.Println(GetVerifyCode())
	}
}
