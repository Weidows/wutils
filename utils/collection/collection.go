package collection

import "github.com/Weidows/Golang/utils/log"

var (
	logger = log.GetLogger()
)

func ForEach[T any | []any](array []T, f func(index int, value T)) {
	for i, v := range array {
		f(i, v)
	}
}

type mapKeys interface {
	int | int32 | int64 | string | float32 | float64
}

func MapToSlice[K mapKeys, V any](m map[K]V) (keys []K, values []V) {
	for i, v := range m {
		keys = append(keys, i)
		values = append(values, v)
	}
	return keys, values
}
