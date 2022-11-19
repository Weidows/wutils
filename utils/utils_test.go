package utils

import (
	"fmt"
	"os"
	"testing"
)

func Test_GetVerifyCode(t *testing.T) {
	for i := 0; i < 15; i++ {
		fmt.Println(GetRandNum(i))
	}
}

func TestGetVerifyCode(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(GetVerifyCode())
	}
}

func TestForEach(t *testing.T) {
	var arr = []int{
		1, 2, 3, 4, 5,
	}
	ForEach(arr, func(i, v int) {
		println(i, v)
	})
}

func TestMapToArray(t *testing.T) {
	keys, values := MapToArray(map[string]int{
		"1": 11,
		"2": 22,
	})
	fmt.Println(keys)
	fmt.Println(values)
}

func TestConditionalEqual(t *testing.T) {
	fmt.Println(ConditionalEqual(true, 1, 2))
	fmt.Println(ConditionalEqual(false, 1, 2))
}

func TestGetMP4Duration(t *testing.T) {
	file, err := os.Open("D:/Scoop/persist/steam/steamapps/workshop/content/431960/2882829381/3840x2160pro_1f460.mp4")
	if err != nil {
		panic(err)
	}
	duration, err := GetMP4Duration(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(duration)
}

func TestMarshal(t *testing.T) {
	data := Marshal(map[string]string{
		"duration": "16",
		"hash":     "b2146aae9e9f807b7b8c17fcc531addedaff6670a37481cf185133be42a31d25",
		"height":   "720",
		"name":     "video.mp4",
		"path":     "/mnt/test_data/Movie/1Marvel/10001",
		"size":     "6901",
		"width":    "1280",
	})
	fmt.Println(data)
}
