/*
 * @?: *********************************************************************
 * @Author: Weidows
 * @Date: 2022-05-25 17:34:28
 * @LastEditors: Weidows
 * @LastEditTime: 2022-05-25 22:38:23
 * @FilePath: \Golang\src\demo\concurrent.go
 * @Description: 简易协程调用/通信
	https://www.runoob.com/go/go-concurrent.html
 * @!: *********************************************************************
*/

package concurrent

import (
	"fmt"
	"time"
)

func goroutine() {
	var say = func(s string) {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			fmt.Println(s)
		}
	}

	go say("hello")
	go say("world")
	// 如果不加sleep的话, main线程太快, 还没say出来就结束了
	time.Sleep(5 * time.Second)
}

func channel() {
	var say = func(s string, c chan string) {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			c <- s
		}
	}
	//c := make(chan string)
	// 缓冲区为2
	c := make(chan string, 2)
	go say("hello", c)
	go say("world", c)

	// 最后因为读空 channel, 导致所有 goroutine 被阻塞(死锁) 而报错退出
	//for i := range c {
	//	fmt.Println(i)
	//}

	count := 10
	for result, isNotEmpty := <-c; result != "" && isNotEmpty; result, isNotEmpty = <-c {
		fmt.Println(result)
		count--
		if count == 0 {
			close(c)
		}
	}
}
