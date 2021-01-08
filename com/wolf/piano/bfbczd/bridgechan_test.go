package bfbczd

import (
	"fmt"
	"testing"
)

func TestBridgeChan(t *testing.T) {
	//创建了10个通道，每个通道都写入一个元素，并将这些通道传递给bridge：
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}

// 返回所有传入bridge的通道中的值
var bridge = func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for { // 2,临时变量
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream: // 若关闭则退出
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range orDone(done, stream) { // 3
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

//2.该循环负责从chanStream中提取通道并将其提供给嵌套循环以供使用。
//3.该循环负责读取已经给出的通道的值，并将这些值重复到valStream中。当前正在循环的流关闭时，我们跳出执行从该通道读取的循环，并继续下一次循环来选择要读取的通道。 这为我们提供了一个不间断的流。
