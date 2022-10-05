package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
)

func main() {
	//	https://blog.csdn.net/guyan0319/article/details/107182205/
	imageFile, _ := os.Open("D:/Pictures/Aimage/w/1e59a367-cd8a-4cb3-9a48-a7faba395861/result.jpg")
	imageConfig, _, _ := image.DecodeConfig(imageFile)
	fmt.Println(imageConfig.Height, imageConfig.Width)

	//img, _, error := image.Decode(imageFile)
	//if error != nil {
	//	fmt.Println(error)
	//}
	//height := img.Bounds().Dy()
	//width := img.Bounds().Dx()
	//fmt.Println(height, width)
}
