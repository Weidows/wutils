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
		logger.Println(i, v)
	})
}

func TestMap2Slice(t *testing.T) {
	keys, values := Map2Slice(map[string]int{
		"1": 11,
		"2": 22,
	})
	logger.Println(keys)
	logger.Println(values)
}

func TestSlice2Map(t *testing.T) {
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
}
