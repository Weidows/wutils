/*
 * @?: *********************************************************************
 * @Author: Weidows
 * @Date: 2022-03-21 16:59:21
 * @LastEditors: Weidows
 * @LastEditTime: 2022-05-25 22:18:09
 * @FilePath: \Golang\src\demo\cli-loader.go
 * @Description: cli-进度条
 * @!: *********************************************************************
 */

package main

import (
	"time"

	"github.com/cheggaaa/pb/v3"
)

func main() {
	count := 1000
	// create and start new bar
	// bar := pb.StartNew(count)

	// start bar from 'default' template
	// bar := pb.Default.Start(count)

	// start bar from 'simple' template
	bar := pb.Simple.Start(count)

	// start bar from 'full' template
	// bar := pb.Full.Start(count)

	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.Finish()
}
