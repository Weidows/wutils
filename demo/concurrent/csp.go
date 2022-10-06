package concurrent

import (
	"fmt"
	"time"
)

var ch = make(chan string)

func receive() {
	// 异步, 阻塞读
	msg := <-ch
	fmt.Println(msg)
}

func send() {
	// 写
	ch <- "Hello,CSP."
}

// channel 并没有读写顺序
func main() {
	go send()
	go receive()

	// 这俩倒过来也没问题
	go receive()
	go send()

	// 这样也没问题
	go receive()
	ch <- "Hello,CSP."

	//* 这样就会出错, 因为没有接收方要阻塞发送方(main), 而阻塞 main 产生了死锁
	ch <- "Hello,CSP."
	go receive()

	// 有 goroutine 的 main 里记得 sleep, 不然还没执行就结束程序了
	time.Sleep(time.Second)
}
