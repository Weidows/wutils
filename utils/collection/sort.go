package collection

import (
	"github.com/Weidows/Golang/utils/cast"
	"github.com/Weidows/Golang/utils/reflect"
	"github.com/gogo/protobuf/sortkeys"
	"sort"
)

// SortKeys
//
// From: https://studygolang.com/articles/10530
func SortKeys[K MapKeys, V any, T interface{ []K }](m map[K]V) T {
	if len(m) == 0 {
		return nil
	}

	keys, _ := MapToSlice(m)
	sorted := SortSlice(keys)
	//res := make([]K, len(keys))
	//for i, v := range sorted.([]K) {
	//	res[i] = reflect2.ValueOf(v).Convert(reflect2.TypeOf(keys[0])).Interface().(K)
	//}
	return sorted
}

// SortableTypes : MapKeys 的超集
type SortableTypes interface {
	int | int8 | int32 | int64 | string | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64 | any
}

func SortSlice[T SortableTypes](slice []T) (sorted []T) {
	reflect.TypeofT(slice, func(a any) {
		switch t := a.(type) {
		case []int:
			sort.Ints(a.([]int))
		case []int8:
			s := cast.Convert[[]int](a)
			sort.Ints(s)
			a = cast.Convert[[]int8](s)
		case []int32:
			sortkeys.Int32s(a.([]int32))
		case []int64:
			sortkeys.Int64s(a.([]int64))
		case []string:
			sort.Strings(a.([]string))
		case []float32:
			sortkeys.Float32s(a.([]float32))
		case []float64:
			sort.Float64s(a.([]float64))
		default:
			logger.Println("Unsupported type")
			_ = t
		}
		sorted = a.([]T)
	})
	return sorted
}
