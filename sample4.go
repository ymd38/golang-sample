package main

import (
	"fmt"
	"time"
)

func main() {
	iCh := make(chan int, 1)
	qCh := make(chan int, 1)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i, ":", <-iCh)
		}
		qCh <- 0
	}()
	send(iCh, qCh)
}

func send(iCh, qCh chan int) {
	x, y := 0, 1
	for {
		select {
		case iCh <- x:
			fmt.Println("send:", x)
			x, y = y, x+y
			time.Sleep(10000 * time.Millisecond)
		case <-qCh:
			return
		}
	}
}
