package json

import (
	"fmt"
	"net/http"
	"testing"
)

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

func TestDecode(t *testing.T) {
	resp, err := http.Get("https://api.btstu.cn/yan/api.php?charset=utf-8&encode=json")
	if err != nil {
		return
	}
	res := Decode[map[string]interface{}](resp.Body)
	fmt.Println(res)
}
