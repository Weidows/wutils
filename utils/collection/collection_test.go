package collection

import (
	"fmt"
	"testing"
)

func TestForEach(t *testing.T) {
	var arr = []int{
		1, 2, 3, 4, 5,
	}
	ForEach(arr, func(i, v int) {
		println(i, v)
	})
}

func TestMapToArray(t *testing.T) {
	keys, values := MapToArray(map[string]int{
		"1": 11,
		"2": 22,
	})
	fmt.Println(keys)
	fmt.Println(values)
}
