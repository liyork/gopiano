package bfbczd

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// 假设我们有一个固定长度为2的队列，并且我们要将10个元素放入队列中。 我们希望一有空间就能放入，所以在队列中有空间时需要立刻通知：
func TestProductConsumer(t *testing.T) {
	c := sync.NewCond(&sync.Mutex{})    //1
	queue := make([]interface{}, 0, 10) //2

	// 生产者和消费者之间用的都是锁，所以同一时间只能有一个进行写或读
	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay) // 每次启动消费者都睡眠1s，给足生产者机会强锁
		c.L.Lock()        //8
		i := queue[0]
		fmt.Println("i:", i)
		queue = queue[1:] //9
		fmt.Println("Removed from queue")
		c.L.Unlock() //10
		c.Signal()   //11
	}

	// 每次运行都上锁，同一时间，最多有两个消费者，肯定有一个能把生产者唤醒
	// 队列长度为2时等待removeFromQueue协程消费后唤醒再继续
	for i := 0; i < 10; i++ {
		c.L.Lock()            //3
		for len(queue) == 2 { //4 ,持续监控，若不为2则继续等待，否则2表示满了
			fmt.Println("queue len", 2)
			c.Wait() //5
		}
		fmt.Println("Adding to queue")
		queue = append(queue, i)
		go removeFromQueue(1 * time.Second) //6
		c.L.Unlock()                        //7
	}
	//1.首先，我们使用一个标准的sync.Mutex作为Locker来创建Cond。
	//2.接下来，我们创建一个长度len为零的切片。 由于我们知道最终会添加10个元素，因此我们将其容量cap设为10。
	//3.在进入关键的部分前调用Lock来锁定c.L。
	//4.在这里我们检查队列的长度，以确认什么时候需要等待。由于removeFromQueue是异步的，for不满足时才会跳出，而if做不到重复判断，这一点很重要。--避免假唤醒
	//5.调用Wait，这将阻塞main goroutine，直到接受到信号。
	//6.这里我们创建一个新的goroutine，它会在1秒后将元素移出队列。
	//7,这里我们退出条件的关键部分，因为我们已经成功加入了一个元素。
	//8,我们再次进入该并发条件下的关键部分，以修改与并发条件判断直接相关的数据。
	//9.在这里，我们移除切片的头部并重新分配给第二个元素，这一步模拟了元素出列。
	//10.我们退出操作关键部分，因为我们已经成功移除了一个元素。
	//11.这里，我们发出信号，通知处于等待状态的goroutine可以进行下一步了。
}

// 为了了解Broadcast是如何使用的，
// 假设我们正在创建一个带有按钮的GUI程序，该程序需要注册任意数量的函数，当点击按钮时运行这些函数。可以使用Cond的Brocast来通知所有已注册函数
// 如果不是wg WaitGroup，我们可以多次调button.Clicked.Broadcast（），并且每次都将运行这三个处理函数。 这是通道难以做到的，也是使用Cond类型的优势之一。
func TestBroadcast(t *testing.T) {
	type Button struct {
		Clicked *sync.Cond //1
	}
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) { //2
		var tempwg sync.WaitGroup
		tempwg.Add(1)
		go func() {
			tempwg.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		tempwg.Wait()
	}

	var wg sync.WaitGroup //3
	wg.Add(3)

	subscribe(button.Clicked, func() { //4
		fmt.Println("Maximizing window.")
		wg.Done()
	})
	subscribe(button.Clicked, func() { //5
		fmt.Println("Displaying annoying dialog box!")
		wg.Done()
	})
	subscribe(button.Clicked, func() { //6
		fmt.Println("Mouse clicked.")
		wg.Done()
	})

	button.Clicked.Broadcast() //7

	wg.Wait()
	//1.我们定义一个Button类型，包含了sync.Cond指针类型的Clicked属性，这是goroutine接收通知的关键条件。
	//2.这里我们定义了一个较为简单的函数，它允许我们注册函数来处理信号。每个注册的函数都在自己的goroutine上运行，并且在该goroutine不会退出，直到接收到通知。
	//3.在这里，我们为按钮点击设置了一个处理程序。 它反过来在Clicked Cond上调用Broad cast以让所有注册函数知道按钮已被点击。
	//4.这里我们创建一个WaitGroup。 这只是为了确保我们的程序在写入标准输出之前不会退出。
	//5.在这里我们注册一个处理函数，模拟点击时最大化窗口。
	//6,在这里我们注册一个处理函数，模拟点击时显示对话框。
	//7.接下来，我们模拟按钮被点击。
}
