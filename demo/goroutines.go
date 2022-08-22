/*
 * @?: *********************************************************************
 * @Author: Weidows
 * @Date: 2022-05-25 17:34:28
 * @LastEditors: Weidows
 * @LastEditTime: 2022-05-25 22:38:23
 * @FilePath: \Golang\src\demo\goroutines.go
 * @Description: 简易协程调用
 * @!: *********************************************************************
 */

package main

import (
	"fmt"
	"time"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		fmt.Println(s)
	}
}

func main() {
	go say("world")
	say("hello")
}
