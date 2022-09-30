package utils

import (
	"fmt"
	"testing"
)

func Test_GetVerifyCode(t *testing.T) {
	for i := 0; i < 15; i++ {
		fmt.Println(GetRandNum(i))
	}
}

func TestGetVerifyCode(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(GetVerifyCode())
	}
}
