package main

import (
	"fmt"
	"time"
)

// 在需要频繁使用相对到期时间相同的定时器的情况下，总是应该尽量复用，而不是重新创建。
//若在定时器到期之前停止了它，那么该定时器的字段C就没有机会缓冲任何元素值了。具体将，若调用Stop结果为true，那么在这之后再去试图从他的C字段中接收
//元素是不会有任何结果的，还会使当前goroutine永远阻塞。因此，在重置定时器之前一定不要再次对他的C字段执行接收操作。
func main() {
	intChan := make(chan int, 1)
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			intChan <- i
		}
		close(intChan)
	}()

	// 每次接收时都使用同一个定时器
	timeout := time.Millisecond * 500
	var timer *time.Timer
	for {
		if timer == nil {
			timer = time.NewTimer(timeout)
		} else {
			timer.Reset(timeout)
		}
		select {
		case e, ok := <-intChan:
			if !ok {
				fmt.Println("End.")
				return
			}
			fmt.Printf("Received: %v\n", e)
		case <-timer.C:
			fmt.Println("Timeout!")
		}
	}
}
