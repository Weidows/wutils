package collection

import "testing"

func TestSortKeys(t *testing.T) {
	logger.Println(SortKeys(map[int]string{
		1: "a",
		3: "b",
		2: "c",
		0: "d",
		5: "e",
	}))
	logger.Println(SortKeys(map[int8]string{
		1: "a",
		3: "b",
		2: "c",
		0: "d",
		5: "e",
	}))
	logger.Println(SortKeys(map[int64]string{
		1: "a",
		3: "b",
		2: "c",
		0: "d",
		5: "e",
	}))
	logger.Println(SortKeys(map[string]string{
		"1": "a",
		"3": "b",
		"2": "c",
		"0": "d",
		"5": "e",
	}))
	logger.Println(SortKeys(map[float32]string{
		1: "a",
		3: "b",
		2: "c",
		0: "d",
		5: "e",
	}))
	logger.Println(SortKeys(map[float64]string{
		1: "a",
		3: "b",
		2: "c",
		0: "d",
		5: "e",
	}))
}
