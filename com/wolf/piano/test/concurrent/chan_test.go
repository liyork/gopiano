package concurrent

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

// channel的作用就是在多线程之间传递数据的
// channel的机制是先进先出，如果你给channel赋值了，那么必须要读取它的值，不然就会造成阻塞，当然这个只对无缓冲的channel有效

// chan的结构是在src/runtime/chan.go 的hchan。我们就看chan.go里面的recv方法
// func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {   这个函数就是<-ch 的时候调用的
// c代表的就是我们使用的这个chan, ep代表的是ch传输出来的数据存储的位置。
// 在从channel中获取数据的时候调用的是 recv(c, sg, ep, func() { unlock(&c.lock) }, 3)  -->typedmemmove,这个函数就是c中的memmove。
// 将原先的数据，直接拷贝到目标内存中。所以这里说明channel是进行值拷贝的(都是值，一般的变量值或者指针值)。

// close没有make的chan会引起panic
//close以后不能再写入，写入会出现panic
//close之后可以读取，无缓冲chan读取返回0值和false，有缓冲chan可以继续读取，返回的都是chan中数据和true，直到读取完所有队列中的数据。
//重复close会引起panic
//只读chan不能close
//不close chan也是可以的，当没有被引用时系统会自动垃圾回收。
//不要在receiver中close(chan)，而要在sender中close(chan)
//多个sender时，要使用辅助chan来确保不会多次close(chan)和close后又再次写入。
//此时最好不要主动close chan而应该在无引用时由系统自动回收，要主动close除非确保所有写入都停止了。

// 不带缓冲的Channel兼具通信和同步两种特性

//无缓冲的信道永远不会存储数据，只负责数据的流通
//从无缓冲信道取数据，必须要有数据流进来才可以，否则当前线阻塞
//数据流入无缓冲信道, 如果没有其他goroutine来拿走这个数据，那么当前线阻塞
//所以，你可以测试下，无论如何，我们测试到的无缓冲信道的大小都是0 (len(channel))
//
//为什么会死锁？非缓冲信道上如果发生了流入无流出，或者流出无流入，也就导致了死锁。或者这样理解 Go启动的所有goroutine里的非缓冲信道一定要一个线里存数据，一个线里取数据，要成对才行
//
//chan像阻塞队列
//
//缓冲这个词意思是，缓冲信道不仅可以流通数据，还可以缓存数据。它是有容量的，存入一个数据的话 , 可以先放在信道里，不必阻塞当前线而等待该数据取走。
//当缓冲信道达到满的状态的时候，就会表现出阻塞了，因为这时再也不能承载更多的数据了
//
//缓冲信道是先进先出的，我们可以把缓冲信道看作为一个线程安全的队列

// go中通过channel进行消息通信，channel函数传参通过引用实现

// channel介绍
//golang社区口号：不要通过共享内存来通信，而应该通过通信来共享内存。 通过消息来传递内存
//golang提供一种基于消息机制而非共享内存的通信模型。消息机制认为每个并发单元都是自包含的独立个体，并且拥有自己的变量，但在不同并发单元间这些变量不共享。每个并发单元的输入和输出只有一种，那就是消息
//channel是golang在语言级提供的goroutine间的通信方式，可以使用channel在两个或多个goroutine之间传递消息
//channel是进程内的通信方式，如果需要跨进程通信，建议使用分布式的方法来解决，比如使用Socket或HTTP等通信协议
//channel是类型相关的，即一个channel只能传递一种类型的值，需要在声明channel时指定。可以认为channel是一种类型安全的管道
//
//var chanName chan ElementType
//var ch chan int            // int类型channel
//var m map[string]chan bool // bool类型channel的map
//
//应该在生产者处关闭channel，而不是消费者处关闭channel，否则容易引起panic。
//
//channel读写语法
//向无缓冲的channel写入数据会导致该goroutine阻塞，直到其他goroutine从这个channel中读取数据
//向带缓冲的且缓冲已满的channel写入数据会导致该goroutine阻塞，直到其他goroutine从这个channel中读取数据
//向带缓冲的且缓冲未满的channel写入数据不会导致该goroutine阻塞
//从无缓冲的channel读出数据，如果channel中无数据，会导致该goroutine阻塞，直到其他goroutine向这个channel中写入数据
//从带缓冲的channel读出数据，如果channel中无数据，会导致该goroutine阻塞，直到其他goroutine向这个channel中写入数据
//从带缓冲的channel读出数据，如果channel中有数据，该goroutine不会阻塞
//总结：无缓冲的channel读写通常都会发生阻塞，带缓冲的channel在channel满时写数据阻塞，在channel空时读数据阻塞
//阻塞同步队列

// range不等到信道关闭是不会结束读取的，会一直阻塞读取，直到close
func TestChanRangeErr(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3

	for v := range ch {
		fmt.Println(v)
	}
}

// 读到信道为空的时候就结束读取:
// 注意:检查信道大小的方法不能在信道存取都正在发生的时候，因为并发写入读取会导致len不准确，只可能用于对[不再写的chan]进行len
func TestChanRangeCorrectWarn(t *testing.T) {
	ch := make(chan int, 3)
	go func() {
		for i := 0; i < 100; i++ {
			ch <- i
		}
	}()

	for v := range ch {
		fmt.Println(v)
		len := len(ch)
		fmt.Println("len:", len)
		if len <= 0 { // 如果现有数据量为0，跳出循环
			break
		}
	}
}

// 显式地关闭信道
func TestChanRangeCorrect(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3

	// 显式地关闭信道
	close(ch)

	// 被关闭的信道会禁止数据流入, 是只读的。可以从关闭的信道中取出剩余数据。
	for v := range ch {
		fmt.Println(v)
		time.Sleep(time.Second)
	}
}

// Closing the nil channel.
func TestCloseNil(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover has err,", err)
		}
	}()

	var ch chan int
	close(ch)
}

// Closing a closed channel.
func TestCloseClosedChan(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover has err,", err)
		}
	}()

	ch := make(chan int, 0)
	close(ch)
	close(ch) // 关闭已关闭的chan
}

// 读写nil的chan都会阻塞
func TestChanNilBlock(t *testing.T) {
	var c chan int = nil
	<-c
	//c <- 1
}

// Sending on a closed channel.
func TestSendClosedChan(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover has err,", err)
		}
	}()

	ch := make(chan int, 0)
	close(ch)
	ch <- 2 // 对close的chan再发送
}

// 带缓冲的channel close，依旧可以从中读出关闭前写入的3个值。第四次读取时，则会返回该channel类型的零值。
// 向这类channel写入操作也会触发panic。
func TestCloseBufChan(t *testing.T) {
	c := make(chan int, 3)
	c <- 15
	c <- 34
	c <- 65
	close(c)
	fmt.Printf("%d\n", <-c)
	fmt.Printf("%d\n", <-c)
	fmt.Printf("%d\n", <-c)
	fmt.Printf("%d\n", <-c)

	// panic
	c <- 1
}

func TestAssignChan(t *testing.T) {
	ch := make(chan int)
	ch1 := ch
	fmt.Println("==,", ch == ch1)
}

// 关于chan关闭时机:
// 作为接受者的goroutine不要关闭channel,否则 发送者未来往这个channel发送数据时会panic
// 如果一个channel有多个发送者，发送者就不要关闭通道或者最后一个发送数据的goroutine来关闭通道
// Last one sender to leave, turns off the lights, which can be controlled by a atomic int
// 对于不再使用的通道不必显示关闭。如果没有goroutine引用这个通道，这个通道就会被垃圾回收
// 注意如果需要把关闭通道作为一个控制信号告知其他goroutine没有更多数据的情况下，需要显示关闭
// 通道和select搭配最佳
// 如果在两个goroutine中需要一个双向沟通，考虑使用两个单独的单向通道。这样两个通道就可以通过调用close的方法来告知对方通讯终止
// 如果通道没有读取者的goroutine的时候，那么发送者就可能会永远阻塞在该通道上。
// 当设计一个goroutine通过channel来提供服务的时候，在某个时刻这个goroutine不再需要了，要仔细考虑goroutine应该怎么结束。否则这个不再使用的goroutine会服务一个不受监管的channel
// select从来不会选择阻塞的case,若多个case都阻塞了，那么当前goroutine就会阻塞在select上。

func convertOnlyReadChan(writeChan chan<- int) {

}

// 只读
func convertOnlyWriteChan(readChan <-chan int) {
	for a := range readChan {
		fmt.Println(a)
	}
	// 编译报错,不允许对只读chan写
	//readChan<-1
}

func testConvert() {
	bichan := make(chan int)
	convertOnlyReadChan(bichan)  // 转成只读
	convertOnlyWriteChan(bichan) // 转成只写

	bichan2 := make(chan int, 2) // capacity
	fmt.Println(bichan2)
}

// 一个chan多个读取只能有一个成功
func TestMultiReceive(t *testing.T) {
	pubchan := make(chan int)

	go func() {
		for {
			a := <-pubchan
			fmt.Println("receive a1 succss, a2:", a)
		}
	}()

	go func() {
		for {
			a := <-pubchan
			fmt.Println("receive a2 succss, a2:", a)
		}
	}()

	for i := 0; i < 10; i++ {
		//time.Sleep(1 * time.Second)
		pubchan <- i
	}
}

// 无缓冲channel则写入/读取阻塞直到另一端有人读/写
func TestNoBufferBlockRead(t *testing.T) {
	noBuffChan := make(chan int)
	go writeNoBuffer(noBuffChan)
	time.Sleep(time.Second * 2)
	noBuffChan <- 2
	fmt.Println("main finish:")
}

func writeNoBuffer(noBuffChan chan int) {
	a := <-noBuffChan
	// 阻塞，直到有人写入
	fmt.Println("a:", a)
}

func TestBufferBlockRead(t *testing.T) {
	buffChan := make(chan int, 2)
	a := <-buffChan // 阻塞，因为缓冲中没有数据
	fmt.Println("a:", a)
}

// 缓冲区满之前都不阻塞
func TestBufferBlockWrite(t *testing.T) {
	buffChan := make(chan int, 2)
	buffChan <- 1
	buffChan <- 2
	buffChan <- 3 // 阻塞，缓冲被写满
	fmt.Println("a:")
}

// 利用close而不是每个发一遍
func TestChannelClose(t *testing.T) {
	chann := make(chan int)
	go processWithChanel(chann)
	fmt.Println("wait processWithChanel")
	go func() {
		i, ok := <-chann
		fmt.Println("finish main1,", i, ok)
	}()
	go func() {
		i, ok := <-chann
		fmt.Println("finish main2,", i, ok)
	}()
	i, ok := <-chann
	fmt.Println("finish main3,", i, ok)
}

var cx = make(chan int, 10)
var ax string

func f() {
	ax = "hello, world"
	cx <- 0
}

// A send on a channel happens before the corresponding receive from that channel completes.
// is guaranteed to print "hello, world".
// The write to ax happens before the send on cx, which happens before the corresponding receive on cx completes,
// which happens before the println.
func TestSendSeqWithBufChan(t *testing.T) {
	go f()
	<-cx
	println(ax)
}

func fWithBufChan() {
	ax = "hello, world"
	//cx <- 0
	// replacing c <- 0 with close(c) yields a program with the same guaranteed behavior.
	close(cx)
}

// The closing of a channel happens before a receive that returns a zero value because the channel is closed.
func TestCloseChanSeqWithBuf(t *testing.T) {
	go fWithBufChan()
	<-cx
	println(ax)
}

var c2 = make(chan int)

//var c2 = make(chan int,10)// 改成这样，还真不行。看来nobufchan，read happen before write
var a2 string

// 执行after fWithNoBufChan和after TestSendSeqWithNoBufChan随机
// 只能说对nobuf chan，<-c2 happen before c2 <- 0，也就是<-c2上的代码比c2 <- 0之后的早，并不能说明<-c2和c2 <- 0之后的代码顺序
func fWithNoBufChan() {
	fmt.Println("before fWithNoBufChan")
	time.Sleep(2 * time.Second)
	a2 = "hello, world"
	<-c2
	fmt.Println("after fWithNoBufChan")
}

// A receive from an unbuffered channel happens before the send on that channel completes.
// is also guaranteed to println "hello, world".
// The write to a2 happens before the receive on c2, which happens before the corresponding send on c2 completes,
// which happens before the println.
// 注意，对于nobuf的chan，读会happen before写
func TestSendSeqWithNoBufChan(t *testing.T) {
	go fWithNoBufChan()
	fmt.Println("before TestSendSeqWithNoBufChan")
	time.Sleep(2 * time.Second)
	c2 <- 0
	fmt.Println("after TestSendSeqWithNoBufChan")
	println("a2:" + a2)
}

// 对于close，任何chan都不能保证。
func TestCloseSeqWithBufChan(t *testing.T) {
	go fWithNoBufChan()
	close(c2)
	println("a2:" + a2)
}

// 看来对于以上happen before一定要保证明白顺序，或者不要有这种强依赖。

var limit = make(chan int, 3)
var work = make([]func(), 0)

// The kth receive on a channel with capacity C happens before the k+Cth send from that channel completes.for buffered channels
// 读kth比k+Cth发送早
// It allows a counting semaphore to be modeled by a buffered channel: the number of items in the channel corresponds
// to the number of active uses, the capacity of the channel corresponds to the maximum number of simultaneous uses,
// sending an item acquires the semaphore, and receiving an item releases the semaphore.
// This is a common idiom for limiting concurrency. 生产消费缓冲限制
// This program starts a goroutine for every entry in the work list, but the goroutines coordinate using the limit
// channel to ensure that at most three are running work functions at a time.
func TestReceiveK(t *testing.T) {
	for _, w := range work {
		go func(w func()) {
			limit <- 1 //取信号
			w()
			<-limit
		}(w)
	}
	select {}
}

func processWithChanel(chann chan int) {
	fmt.Println("business")
	time.Sleep(time.Second * 3)
	// there is error
	close(chann)
}

// select操作，相对于range能处理多个chan
//golang中的select关键字用于处理异步IO，可以与channel配合使用
//golang中的select的用法与switch语言非常类似，不同的是select每个case语句里必须是一个IO操作
//select会一直等待等到某个case语句完成才结束
//
//select {
//case <-chan1:
//    // 如果chan1成功读到数据，则进行该case处理语句
//case chan2 <- 1:
//    // 如果成功向chan2写入数据，则进行该case处理语句
//default:
//    // 如果上面都没有成功，则进入default处理流程
//}

func TestChanTimeout(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		time.Sleep(122 * time.Second)
		ch <- 1
	}()

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	select {
	case <-ch:
		fmt.Println("read from ch")
	case <-timeout:
		fmt.Println("read from ch timeout")
	}
}

func TestChanTimeAfter(t *testing.T) {
	ch := make(chan int, 1)

	var count = 0
	for {
		if count == 3 {
			fmt.Println("3 times")
			return
		}
		select {
		case <-ch:
			fmt.Println("read from ch")
		case <-time.After(2 * time.Second):
			fmt.Println("time out")
		}
		count++
	}
}

func TestChanStop(t *testing.T) {
	ch := make(chan int, 0)
	stop := make(chan int, 0)
	callback := make(chan int, 0)
	go func() {
		for {
			select {
			case <-ch:
				fmt.Println("read from ch")
			case <-stop:
				fmt.Println("read from stop chan")
				callback <- 2
				break
			}
			time.Sleep(2 * time.Second)
		}
	}()

	for i := 0; i < 3; i++ {
		ch <- i
	}

	stop <- 1
	fmt.Println(<-callback)
}

// 只要一个结果即可
func Query(conns []string, query string) string {
	ch := make(chan string)
	for _, conn := range conns {
		go func(c string) {
			select {
			case ch <- DoQuery(query):
			default: // If the send cannot go through immediately the default case will be selected.
			}
		}(conn)
	}
	return <-ch
}

func DoQuery(s string) string {
	return ""
}

type chanwrapper struct {
	ch chan int
}

// chan通过make出来是指针
func TestChanAssign(t *testing.T) {
	c := chanwrapper{ch: make(chan int)}
	go func() {
		copy := c.ch
		time.Sleep(1 * time.Second)
		//copy <- 1
		fmt.Println(copy == c.ch)
		c.ch = make(chan int) // 先赋值
		time.Sleep(1 * time.Second)
		close(copy) //再关闭老的
	}()
	fmt.Println("wait processWithChanel")
	go func() {
		ints := c.ch // 指针拷贝
		time.Sleep(3 * time.Second)
		fmt.Println("ints == c.ch==>", ints == c.ch) //false，因为上面重新make了
		i, ok := <-ints
		fmt.Println("finish main1,", i, ok)
	}()

	time.Sleep(5 * time.Second)
}

var wg sync.WaitGroup

// 生产者
func Send(shareChan chan int) {
	x := 0
	defer func() {
		if err := recover(); err != nil && err.(runtime.Error).Error() == "send on closed channel" {
			fmt.Println("Send recover has err,", err)
			fmt.Println("最后一次产生的数据：", x)
		} else {
			close(shareChan) //关闭的目的：不在发送数据
		}
		wg.Done()
	}()
	for i := 0; i < 10; i++ {
		x++
		shareChan <- x
	}
}

// 消费者,一直读取shareChan，第三次panic
func Receive(shareChan chan int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Receive recover err:,", err)
			close(shareChan) //不要让生产者继续发送数据
			i, ok := <-shareChan
			fmt.Println("Receive recover receive msg:,", i, ok) //关闭的流再读也是0
		}
		// 通知main完成
		wg.Done()
	}()
	for x := range shareChan {
		time.Sleep(time.Second)
		fmt.Println("read msg:", x)
		if x == 3 {
			panic("发生意外的错误")
		}
	}
	//不会执行，因为上面发生panic,之后被recover然后执行完recover就完了
	fmt.Println("Receive任务结束")
}

// 消费者意外时，关闭close不让provider继续产生无用数据了。
func TestChanProcessErr(t *testing.T) {
	shareChan := make(chan int)
	wg.Add(2)
	go Send(shareChan)
	go Receive(shareChan)
	wg.Wait()
	fmt.Println("main任务完成")
	_, ok := <-shareChan
	fmt.Println("main read from shareChan,", ok)
}

type TestChan struct {
	msgChan  chan int
	stopChan chan struct{}
}

func (t *TestChan) sendMsg(msg int) bool {
	//这个是为了避免stopChan已经close但是下面的第二个select多次随机执行t.msgChan <- msg
	select {
	case <-t.stopChan:
		return false
	default:
	}

	select {
	case <-t.stopChan:
		return false
	case t.msgChan <- msg: //当stopChan close时最多写入一次(随机)，如果不能写则阻塞，但是<-t.stopChan会被激活于是return
	}
	return true
	//我们不主动关闭t.msgChan，在无引用时会自动释放
}
func (t *TestChan) processMsg() {
	for {
		select {
		case <-t.stopChan:
			return
		case msg := <-t.msgChan:
			fmt.Println("processMsg ", msg)
		}
	}
}

// processMsg在这种情况下性能应该更好，因为只需要select一个chan状态，
// 但是这种用法需要自己手动close(t.msgChan)才能退出循环，而且要保证所有的addMsg循环都退出了
// func (t *TestChan) processMsg() {
//  for {
//      select {
//      case msg, ok := <-t.msgChan:
//          if !ok {
//              fmt.Println("processMsg !ok")

//              return
//          }
//          fmt.Println(msg)
//      }
//  }
// }

func waitAndStop(t *TestChan) {
	func() {
		time.Sleep(time.Second * 5)
		fmt.Println("close stopchan")
		close(t.stopChan)
	}()
}

func makechan() *TestChan {
	t := &TestChan{}
	t.msgChan = make(chan int)
	t.stopChan = make(chan struct{})
	return t
}

func processer(t *TestChan) {
	go func() {
		t.processMsg()
		fmt.Println("stop processMsg")
	}()
}

func sender(group *sync.WaitGroup, t *TestChan) {
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(msg int) {
			for {
				if !t.sendMsg(msg) {
					fmt.Println("sendMsg stop", msg)
					wg.Done()
					return
				}
				time.Sleep(time.Second * 1)
			}
		}(i)
	}
}

// 2写1读测试
func TestCloseChan(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	testChan := makechan()
	sender(&wg, testChan)
	processer(testChan)
	waitAndStop(testChan)

	wg.Done()
	wg.Wait()
	fmt.Println("exit")
}

// close掉，chan则退出
func TestForChan(t *testing.T) {
	group := sync.WaitGroup{}
	group.Add(2)
	intchan := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			intchan <- i
			time.Sleep(500 * time.Millisecond)
		}
		close(intchan)
		group.Done()
	}()

	go func() {
		for a := range intchan {
			fmt.Println("read: ", a)
		}
		group.Done()
	}()

	group.Wait()
}

func TestTimeoutChan(t *testing.T) {
	ch := make(chan int)
	select {
	case <-time.After(time.Second * 2):
		println("opt channel timeout")
	case i := <-ch:
		println(i)
	case ch <- 1:
		println("write succ")
	}
}

func waiter(i int, block, done chan struct{}) {
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	fmt.Println(i, "waiting...")
	<-block
	fmt.Println(i, "done!")
	done <- struct{}{}
}

// 根据《Go memory model》中关于close channel与recv from channel的order的定义：
// The closing of a channel happens before a receive that returns a zero value because the channel is closed.
func TestUseForBroadcast(t *testing.T) {
	block, done := make(chan struct{}), make(chan struct{})
	for i := 0; i < 4; i++ {
		go waiter(i, block, done)
	}
	time.Sleep(5 * time.Second)
	close(block)
	for i := 0; i < 4; i++ {
		<-done
	}
}

type Work struct {
	Job string
}

func makeWork(workChan chan Work) {
	for {
		time.Sleep(500 * time.Millisecond)
		workChan <- Work{Job: "job"}
	}
}

func worker(workerIndex int, workChan chan Work, quitChan chan struct{}) {
	for {
		select {
		case w := <-workChan:
			if quitChan == nil {
				fmt.Println(w.Job + "stopped")
				fmt.Println("worker", workerIndex, "refused", w)
				break
			}
			fmt.Println("do", w.Job)
			fmt.Println("worker", workerIndex, "processed", w)
		case <-quitChan:
			fmt.Println("worker", workerIndex, "quitting")
			quitChan = nil
		}
	}
}

// 对于值为nil的channel的case进行接收或发送操作时，会永久block
func TestUseNilChan(t *testing.T) {
	works := make(chan Work)
	quit := make(chan struct{})
	go makeWork(works)
	for i := 0; i < 2; i++ {
		go worker(i, works, quit)
	}
	time.Sleep(5 * time.Second)
	quit <- struct{}{}
}

// 在 select 中动态的打开和关闭 case 语句块
func TestNilChan(t *testing.T) {
	var c1, c2 = make(chan int), make(chan int)
	//var c1, c2 chan int//改成这样,都是nil，那么select会一直阻塞，
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("go 1 awake")
		c1 <- 5
		close(c1)
	}()
	go func() {
		time.Sleep(time.Second * 7)
		fmt.Println("go 2 awake")
		c2 <- 7
		close(c2)
	}()

	fmt.Println("prpare for select")
	for {
		select {
		case x, ok := <-c1:
			if !ok {
				c1 = nil // closed后，设定nil，让select不在对其watch
				fmt.Println("set c1 nil")
			} else {
				fmt.Println(x)
			}
		case x, ok := <-c2:
			if !ok {
				c2 = nil
				fmt.Println("set c1 nil")
			} else {
				fmt.Println(x)
			}
		}
		if c1 == nil && c2 == nil {
			break
		}
	}
	fmt.Println("over")
}

// 读取closed chan返回变量的零值
//判断channel关闭
//在读取的时候使用多重返回值来判断一个channel是否已经被关闭
func TestReadFromClosedChan(t *testing.T) {
	ch := make(chan int, 0)

	close(ch)
	a, ok := <-ch
	fmt.Println(a, ok)
}

// 若是读线程退出了，那么写线程会一直阻塞而不能退出
// 要不就让大家都感知到stop事件，要不就是写事件添加一个超时
func TestBufChanTime(t *testing.T) {
	ch := make(chan int, 2)
	stop := make(chan int)
	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("write receive stop return")
				return
			case ch <- 1:
				fmt.Println("write")
			}

			time.Sleep(time.Second)
		}
	}()
	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("read receive stop return")
				return
			case ch <- 1:
				fmt.Println("read")
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	time.Sleep(5 * time.Second)
	close(stop)

	time.Sleep(55555 * time.Second)
}

// 误解一：有缓存channel是顺序的???这个放入感觉应该是顺序的啊，是不是由于多协成拉取进而打印导致的顺序问题?
// 多次运行会产生不顺序问题
// 1
//0
//2
//3
//4
//5
//6
//7
//8
//9
func TestBufChanIsNoSeq(t *testing.T) {
	cache := make(chan int, 4)
	// 单协成写入10次
	go func() {
		for i := 0; i < 10; i++ {
			cache <- i
		}
	}()
	// 3协成进行读取
	go getCache(cache)
	go getCache(cache)
	go getCache(cache)
	time.Sleep(3 * time.Second)
}

func getCache(cache <-chan int) {
	for {
		select {
		case i := <-cache:
			println(i)
			time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
		}
	}
}

// 误解二：channel缓存的大小就是并发度,多个协成可以从chan中并发同时读取
// 想要并发度，需要不仅chan带有缓冲还要有多协成
func TestBufChanConcurrentCount(t *testing.T) {
	var wg = sync.WaitGroup{}
	wg.Add(2)
	bf := make(chan string, 64)

	go insert(bf, wg)
	// 错误示范
	//go get(bf, wg)
	// 正确操作
	for i := 0; i < 64; i++ {
		go get1(bf, wg)
	}

	wg.Wait()
}

// 不断写入
func insert(bf chan string, wg sync.WaitGroup) {
	str := "模拟无关数据---考生在答题前请认真阅读本注意事项及各题答题要求考生在答题前请认真阅读本注意事项及各题答题要求考生在答题前请认真阅读本注意事项及各题答题要求考生在答题前请认真阅读本注意事项及各题答题要求考生在答题前请认真阅读本注意事项及各题答题要求"
	for i := 0; i < 10000000; i++ {
		bf <- fmt.Sprintf("%s%d", str, i)
	}
	wg.Done()
}

func sprint(str string) {
	fmt.Println("str:", str)
	time.Sleep(1000 * time.Millisecond)
}

// 不断开协成进行读取bf中数据
// 执行就会发现内存一直在增长。 因为get()函数中启动的goroutine会越来越多，
// 因为get()协成很多，导致每读取一个数据，insert()就会往channel插入一条数据，此时并发度就不是64了
func get(bf chan string, wg sync.WaitGroup) {
	for {
		go func() {
			select {
			case str := <-bf:
				sprint(str)
			case <-time.After(3 * time.Second):
				wg.Done()
			}
		}()
	}
}

// 仅一协成for执行读
func get1(bf chan string, wg sync.WaitGroup) {
	for {
		select {
		case str := <-bf:
			sprint(str)
		case <-time.After(3 * time.Second):
			wg.Done()
		}
	}
}

// 协程泄露问题，对于channel的操作可能引发协程泄露，即协程不能被系统垃圾回收，阻塞于chan
// 由于go存在垃圾回收机制，所以无需显示关闭channel，即使是有缓存的channel且channel不为空，也会被垃圾回收，只要保证不再被引用，称为垃圾
// 一个是routine泄漏，即被阻塞不能正常执行完，一个是chan回收，没有被使用则可以回收
func TestGoroutineLeak(t *testing.T) {
	q := make(chan string)
	go func() { q <- "hello" }()
	go func() { q <- "world" }()
	go func() { q <- "china" }()
	<-q
	fmt.Println("111")
}

func MyRoutineFunc(ch chan int) {
	// 函数处理
	ch <- 1

	fmt.Println("MyRoutineFunc process finished.")
}

// 主函数等待所有goroutine完成后返回
//我们已经知道golang程序从main()函数开始执行，当main()函数返回时，程序结束且不等待其他goroutine结束。如果main函数使用time.Sleep方式阻塞等待所有goroutine返回，那么这个休眠时间势必无法控制精确。通过使用channel可以很好解决这个问题
func TestWaitForAll(t *testing.T) {
	chs := make([]chan int, 10)

	for i := 0; i < 10; i++ {
		chs[i] = make(chan int)
		go MyRoutineFunc(chs[i])
	}

	for _, ch := range chs {
		<-ch
	}

	fmt.Println("All goroutine finished.")
}

//WorkChanOut 工作流
func WorkChanOut(i int, sd chan int, stopC chan<- bool) {
	for {
		sum := <-sd
		if sum == 50 {
			stopC <- true
			return
		}
		sum++
		log.Println(i, ":", sum)
		sd <- sum
	}
}

// 三个线程的开始排序是不同的，接下来的输出都是以这种顺序输出(012)
func TestRoutineChanSeq(t *testing.T) {
	sd := make(chan int)
	stopC := make(chan bool)
	for i := 0; i < 3; i++ {
		go WorkChanOut(i, sd, stopC)
	}
	sd <- 0
	log.Println("stop", <-stopC)
}
