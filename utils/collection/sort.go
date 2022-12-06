package collection

import (
	cast2 "github.com/Weidows/Golang/utils/cast"
	"github.com/Weidows/Golang/utils/reflect"
	"github.com/cuigh/auxo/util/cast"
	"sort"
)

func SortByRow() {

}

func SortByColumn() {

}

// SortByKeys
//
// From: https://studygolang.com/articles/10530
func SortByKeys[K MapKeys, V any, T interface{ []K }](m map[K]V) T {
	if len(m) == 0 {
		return nil
	}

	keys, _ := MapToSlice(m)
	var sorted any
	reflect.Typeof(keys, func(a any) {
		switch t := a.(type) {
		case []int:
			sort.Ints(a.([]int))
			sorted = a
		case []int32:
			s := cast2.ToIntSlice(a)
			sort.Ints(s)
			sorted = cast.ToInt32Slice(s)
		case []int64:
			s := cast2.ToIntSlice(a)
			sort.Ints(s)
			sorted = cast.ToInt64Slice(s)
		case []string:
			sort.Strings(a.([]string))
			sorted = a
		case []float32, []float64:
			f := cast2.ToFloat64Slice(a)
			sort.Float64s(f)
			sorted = f
		default:
			_ = t
		}
	})

	//res := make([]K, len(keys))
	//for i, v := range sorted.([]K) {
	//	res[i] = reflect2.ValueOf(v).Convert(reflect2.TypeOf(keys[0])).Interface().(K)
	//}
	return sorted.([]K)
}
