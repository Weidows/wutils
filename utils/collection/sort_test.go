package collection

import "testing"

func TestSortByKeys(t *testing.T) {
	m := map[int64]string{
		1: "a",
		3: "b",
		2: "c",
		0: "d",
		5: "e",
	}
	sorted := SortByKeys[int64](m)
	logger.Println(sorted)
}
