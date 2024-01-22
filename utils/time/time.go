package time

import (
	"github.com/Weidows/wutils/utils/cast"
	"github.com/Weidows/wutils/utils/log"
	"time"
)

var (
	logger = log.GetLogger()
)

// WithTimeOut 超时退出返回 nil, 注意 func 返回值不能是 any
//
// From: https://geektutu.com/post/hpg-timeout-goroutine.html
//
//	WithTimeOut(800*time.Millisecond, func () int {
//		time.Sleep(time.Millisecond)
//		return 1
//	}) // 1
//
//	WithTimeOut(800*time.Millisecond, func () string {
//		time.Sleep(time.Second)
//		return "2"
//	}) // nil
func WithTimeOut[T any](timeout time.Duration, fn func() T) T {
	done := make(chan T, 1)

	go func() {
		// fmt.Println(unsafe.Sizeof(struct{}{}))  // 0
		//done <- struct{}{}
		done <- fn()
	}()

	// 阻塞等待
	select {
	case d := <-done:
		return d
	case <-time.After(timeout):
		//fmt.Println("timeout!!!")
		return cast.EmptyT[T]()
	}
}

func TimeCosts(fn func()) time.Duration {
	t := time.Now()
	fn()
	return time.Now().Sub(t)
}
