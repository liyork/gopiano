package main

import (
	"fmt"
	"time"
)

type Counter struct {
	count int
}

func (counter *Counter) String() string {
	return fmt.Sprintf("{countVal:%d}", counter.count)
}

func main() {
	// 结构体不是引用类型，传值就是拷贝值本身
	//var mapChan = make(chan map[string]Counter, 1)
	// 指针是引用类型，传值也是拷贝值本身，不过是指针的拷贝，最后指向的还是一个变量的内存地址
	var mapChan = make(chan map[string]*Counter, 1)

	syncChan := make(chan struct{}, 2)
	go func() { // 接收
		for {
			if elem, ok := <-mapChan; ok {
				counter := elem["count"]
				counter.count++
			} else {
				break
			}
		}
		fmt.Println("Stopped. [receiver]")
		syncChan <- struct{}{}
	}()

	go func() { // 发送
		// 要使用Counter就会打印map中字段+值
		//countMap := map[string]Counter{"count": Counter{},}
		// 要使用*Counter，就会打印指针的地址，需要重写String
		countMap := map[string]*Counter{"count": &Counter{}}
		for i := 0; i < 5; i++ {
			mapChan <- countMap
			time.Sleep(time.Millisecond)
			fmt.Printf("The count map: %v. [sender]\n", countMap)
		}
		close(mapChan)
		syncChan <- struct{}{}
	}()

	<-syncChan
	<-syncChan
}
