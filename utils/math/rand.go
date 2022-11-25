package math

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	fmt.Println("math")
}

func GetRandNum(digit int) (res string) {
	if digit < 1 {
		return ""
	}
	if digit < 10 {
		// int32 封顶 2^9, 2^10 会溢出
		num := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(int32(Pow10(digit)))
		res = fmt.Sprintf("%0"+strconv.Itoa(digit)+"v", num)
	} else {
		num := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(int32(Pow10(9)))
		res = fmt.Sprintf("%09v", num) + GetRandNum(digit-9)
	}
	return
}

func GetVerifyCode() string {
	return GetRandNum(6)
}
