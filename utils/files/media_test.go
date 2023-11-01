package files

import (
	"fmt"
	"os"
	"testing"
)

func TestGetMP4Duration(t *testing.T) {
	file, err := os.Open("./2e3ce48952af857ccbecb2e8f7ff52c6.mp4")
	if err != nil {
		panic(err)
	}
	duration, err := GetMP4Duration(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(duration)
}
