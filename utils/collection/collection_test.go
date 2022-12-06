package collection

import (
	"testing"
)

func TestForEach(t *testing.T) {
	var arr = []int{
		1, 2, 3, 4, 5,
	}
	ForEach(arr, func(i, v int) {
		logger.Println(i, v)
	})
}

func TestMapToSlice(t *testing.T) {
	keys, values := MapToSlice(map[string]int{
		"1": 11,
		"2": 22,
	})
	logger.Println(keys)
	logger.Println(values)
}
