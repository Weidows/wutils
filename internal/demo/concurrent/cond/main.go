package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	lock                       *sync.Mutex
	consumerCond, producerCond *sync.Cond
	condition                  int
)

func init() {
	lock = new(sync.Mutex)
	consumerCond = sync.NewCond(lock)
	producerCond = sync.NewCond(lock)
	condition = 0
}

// https://cyent.github.io/golang/goroutine/sync_cond/
func main() {
	for i := 0; i < 9; i++ {
		go consumer()
	}

	go producer()
	producer()
}

func consumer() {
	for {
		time.Sleep(time.Second)
		consumerCond.L.Lock()
		for condition == 0 {
			consumerCond.Wait()
		}
		condition--
		fmt.Printf("Consumer: %d\n", condition)
		producerCond.Signal()
		consumerCond.L.Unlock()
	}
}

func producer() {
	for {
		time.Sleep(time.Second)
		producerCond.L.Lock()
		for condition > 10 {
			producerCond.Wait()
		}
		condition += 5
		fmt.Printf("Producer: %d <-----\n", condition)
		consumerCond.Broadcast()
		producerCond.L.Unlock()
	}
}
