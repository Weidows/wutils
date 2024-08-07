package math

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GetRandNum(digit int) string {
	if digit < 1 {
		return ""
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sb := strings.Builder{}
	for sb.Len() < digit {
		sb.WriteString(fmt.Sprint(r.Int63()))
	}
	return sb.String()[:digit]
}

func GetVerifyCode() string {
	return GetRandNum(6)
}
