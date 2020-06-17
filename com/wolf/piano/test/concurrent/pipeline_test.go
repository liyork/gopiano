package concurrent

import (
	"fmt"
	"testing"
	"time"
)

// Golang中并发核心思路是关注数据流动。数据流动的过程交给channel，数据处理的每个环节都交给goroutine
// 流水线由多个阶段组成，每个阶段之间通过channel连接，每个节点可以由多个同时运行的goroutine组成。

// 谁产生chan谁负责收尾售后?
func producer(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer logClose(out)      // 都发送并接收完则关闭
		for _, n := range nums { // 只要有数据就执行循环体
			out <- n // 阻塞,直到有消费
		}
		fmt.Println("producer finish")
	}()
	return out
}

func logClose(out chan int) {
	fmt.Println("logClose finish")
	close(out)
}

func midHandler(inCh <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range inCh {
			out <- n * n
		}
	}()

	return out
}

func consumer(ch <-chan int) {
	for ret := range ch {
		fmt.Printf("%3d", ret)
	}
	fmt.Println()
}

// time go run pipeline_test.go
func TestPipelineBase(t *testing.T) {
	in := producer(1, 2, 3, 4)
	time.Sleep(5 * time.Second)
	ch := midHandler(in)
	consumer(ch)
	fmt.Println()
}
