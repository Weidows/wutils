package collection

import "github.com/Weidows/wutils/utils/log"

var (
	logger = log.GetLogger()
)

func ForEach[T any | []any](array []T, f func(index int, value T)) {
	for i, v := range array {
		f(i, v)
	}
}

type MapKeys interface {
	int | int8 | int32 | int64 | string | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}

func Map2Slice[K MapKeys, V any](m map[K]V) (keys []K, values []V) {
	for i, v := range m {
		keys = append(keys, i)
		values = append(values, v)
	}
	return keys, values
}

// KeySelector 是一个函数类型，用于选择切片元素作为映射的键
type KeySelector[K comparable, V any] func(V) K

/*
// 示例1：将整数切片转换为映射，键为整数值
intSlice := []int{1, 2, 3, 4}

	intMap := Slice2Map(intSlice, func(item int) int {
		return item
	})

fmt.Println(intMap)

// 示例2：将字符串切片转换为映射，键为字符串本身
strSlice := []string{"apple", "banana", "cherry"}

	strMap := Slice2Map(strSlice, func(item string) string {
		return item
	})

fmt.Println(strMap)

// 示例3：将结构体切片转换为映射，键为结构体的某个字段

	type Person struct {
		ID   int
		Name string
	}

	personSlice := []Person{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}

	personMap := Slice2Map(personSlice, func(item Person) int {
		return item.ID
	})

fmt.Println(personMap)
*/
func Slice2Map[K comparable, V any](slice []V, keySelector KeySelector[K, V]) map[K]V {
	result := make(map[K]V)
	for _, item := range slice {
		key := keySelector(item)
		result[key] = item
	}
	return result
}
