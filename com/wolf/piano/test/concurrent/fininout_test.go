package concurrent

import (
	"sync"
	"testing"
)

// FAN-IN和FAN-OUT模式
//FAN-OUT模式：多个goroutine从同一个通道读取数据，直到该通道关闭。OUT是一种张开的模式，所以又被称为扇出，可以用来分发任务。
//FAN-IN模式：1个goroutine从多个通道读取数据，直到这些通道关闭。IN是一种收敛的模式，所以又被称为扇入，用来收集处理的结果。
//FAN模式能够更好的利用CPU，提供更好的并发，提高Golang程序的并发性能。
//在协程比较费时时，FAN模式可以减少程序运行时间，同样的时间，可以处理更多的数据，吞吐量。
//FAN模式不一定能提升效率，降低程序运行时间。
//FAN-IN，squre函数很快就完成，当前程序的瓶颈在merge函数它把3个数据写入到1个通道的时候出现了瓶颈(一个chan有3个写一个读)，适当使用带缓冲通道可以提高程序性能

func square(inCh <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range inCh {
			out <- n * n
		}
	}()
	return out
}

func collectFanin(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	collect := func(in <-chan int) {
		defer wg.Done()
		for n := range in {
			out <- n
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go collect(c)
	}

	go func() {
		// 等到consumer消费完后，collect内的wg.Done执行然后这里进行释放
		wg.Wait()
		close(out)
	}()

	return out
}

// time go run fininout_test.go

func TestFaninout(t *testing.T) {
	in := producer(1, 2, 3, 4)
	c1, c2, c3 := fanout(in)
	collectData := collectFanin(c1, c2, c3)
	consumer(collectData)
}

// 相同逻辑多goroutine执行
func fanout(in <-chan int) (<-chan int, <-chan int, <-chan int) {
	c1 := square(in)
	c2 := square(in)
	c3 := square(in)
	return c1, c2, c3
}
