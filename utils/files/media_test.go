package files

import (
	"fmt"
	"os"
	"testing"
)

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
