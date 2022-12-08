package cast

import (
	"testing"
)

func TestEmptyT(t *testing.T) {
	logger.Println(EmptyT[string]())
	logger.Println(EmptyT[int]())
	logger.Println(EmptyT[bool]())
	logger.Println(EmptyT[[]int]())
	logger.Println(EmptyT[chan string]())
	logger.Println(EmptyT[map[string]interface{}]())
	type a func()
	logger.Println(EmptyT[a]())
}
