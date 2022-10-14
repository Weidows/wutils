package utils

func ForEach[T any](array []T, f func(index int, value T)) {
	for i, v := range array {
		f(i, v)
	}
}

func MapToArray[K, V int | int32 | int64 | string | float32 | float64](m map[K]V) (keys []K, values []V) {
	for i, v := range m {
		keys = append(keys, i)
		values = append(values, v)
	}
	return keys, values
}
