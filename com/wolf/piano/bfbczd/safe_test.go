package bfbczd

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

//可以看到，loopData函数和对handleData通道的循环都使用了整数切片data，但只有loopData对其进行了直接访问
func TestImmutable(t *testing.T) {
	data := make([]int, 4)

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i] // 数据进通道时是份拷贝
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

// 词法约束涉及使用词法作用域仅公开用于多个并发进程的正确数据和并发原语。 这使得做错事情变得不可能
func TestScope(t *testing.T) {
	chanOwner := func() <-chan int {
		results := make(chan int, 5) //1
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ { // 6元素
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) { //3
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}

	results := chanOwner() //2
	consumer(results)
	//1.这里我们在chanOwner函数的词法范围内实例化通道。这将导致通道的写入操作范围被限制在它下面定义的闭包中。 换句话说，它限制了这个通道的写入使用范围，以防止其他goroutine写入它。  --对外就是只能读
	//2.在这里，我们接受到一个只读通道，我们将它传递给消费者，消费者只能从中读取信息。
	//3.这里我们收到一个int通道的只读副本。通过声明该函数的唯一用法是读取访问，我们将通道用法限制为只读。
	//这样的设计方式就可以把通道的读取写入限制在一定的范围内
}

//能并行尽量不并发
func TestConstraint(t *testing.T) {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go printData(&wg, data[:3]) // 1
	go printData(&wg, data[3:]) // 2

	wg.Wait()
	//1.这里我们传入包含前三个字节的data切片。
	//2.这里我们传入包含剩余三个字节的data切片。
}
