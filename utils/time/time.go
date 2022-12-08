package time

import (
	"time"
)

// WithTimeOut 超时退出
//
// var res string
//
//	WithTimeOut(time.Duration(800*time.Millisecond), func() {
//		time.Sleep(time.Millisecond)
//		res += "1"
//
//		time.Sleep(time.Second)
//		res += "2"
//	})
//
// println(res) // 1
func WithTimeOut(timeOut time.Duration, fn func()) {
	done := make(chan any, 1)

	go func() {
		fn()
		done <- nil
	}()

	// 阻塞等待
	select {
	case _ = <-done:
		return
	case <-time.After(timeOut):
		//fmt.Println("timeout!!!")
		return
	}
}
