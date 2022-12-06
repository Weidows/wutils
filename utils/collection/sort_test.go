package collection

import "testing"

func TestSortByKeys(t *testing.T) {
	m := map[int]string{
		1: "a",
		3: "b",
		2: "c",
		0: "d",
		5: "e",
	}
	sorted := SortByKeys[int](m)
	logger.Println(sorted)
}
