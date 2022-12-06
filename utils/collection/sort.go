package collection

import (
	"github.com/Weidows/Golang/utils/reflect"
	"sort"
)

func SortByRow() {

}

func SortByColumn() {

}

// SortByKeys
//
// From: https://studygolang.com/articles/10530
func SortByKeys[K mapKeys, V any](m map[K]V) []K {
	if len(m) == 0 {
		return nil
	}

	keys, _ := MapToSlice(m)
	reflect.Typeof(keys, func(a any) {
		switch t := a.(type) {
		case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64:
			sort.Ints(a.([]int))
		case []string:
			sort.Strings(a.([]string))
		case []float32, []float64:
			sort.Float64s(a.([]float64))
		default:
			_ = t
		}
	})
	return keys
}
