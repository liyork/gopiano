A channel provides a mechanism for concurrently executing functions to communicate by sending and receivingvalues of a specified element type. The value of an uninitialized channel is nil.
未初始化的chan的值是nil(这一点可以看出chan是interface类型，只是内建在编译器内)。

ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType .
chan的类型可以是chan、chan<-、<-chan任意一种+数据类型(如int或者自定义类型)

The optional <- operator specifies the channel direction, send or receive. If no direction is given, the channel isbidirectional. A channel may be constrained only to send or only to receive by conversion or assignment.

The capacity, in number of elements, sets the size of the buffer in the channel. If the capacity is zero or absent, the channel is unbuffered and communication succeeds only when both a sender and receiver are ready. Otherwise, the channel is buffered and communication succeeds without blocking if the buffer is not full (sends) or not empty (receives). A nil channel is never ready for communication.
1. chan在接收和发送会阻塞，阻塞条件是接收时缓冲空或者发送时缓冲满；
2. 如果没有缓冲，接收者和发送者需要同时就绪才会通信，否则调用者就会阻塞，就是接收者调用接收(<-chan)同时发送者调用发送(chan<-)那一刻。
我们常常写测试程序的时候在main函数中创建了一个无缓冲的chan，然后立刻发送一个数据，后面再创建协程接收数据，main函数就会阻塞造成死锁。这是为什么呢？
因为无缓冲chan在双方都就绪后才能通信，否则就会阻塞调用者，所以要先创建协程接收数据，然后再main函数中发送一个数据。
3. 没有被初始化的chan在调用发送或者接收的时候会被阻塞


A channel may be closed with the built-in function close. The multi-valued assignment form of the receive operator reports whether a received value was sent before the channel was closed.
<-chan是一个两个返回值的函数，第一个返回值是指定类型的对象，第二个返回值就是是否接收到了数据，如果第二个返回值是false，说明chan已经关闭了
当chan缓冲中还有一些数据时，关闭chan（调用内建函数close）后，接收者不会立刻收到chan关闭信号(就是第二个返回值为false)，而是等缓冲中所有的数据全部被读取后接收者才会收到chan关闭信号

A single channel may be used in send statements, receive operations, and calls to the built-in functions cap andlen by any number of goroutines without further synchronization. Channels act as first-in-first-out queues. For example, if one goroutine sends values on a channel and a second goroutine receives them, the values are received in the order sent.


go/src/runtime/chan.go

type hchan struct {
	qcount   uint           // 队列中的元素总量
	dataqsiz uint           // 缓冲大小，=make(chan T, x)中的x
	buf      unsafe.Pointer // 缓冲区地址
	elemsize uint16         // 元素大小，单位为字节
	closed   uint32         // chan关闭标记
	elemtype *_type         // 元素类型
	sendx    uint           // 待发送元素在缓冲器中的索引
	recvx    uint           // 待接收元素在缓冲器中的索引
	recvq    waitq          // 接收等待队列，用于阻塞接收协程
	sendq    waitq          // 发送等待队列，用于阻塞发送协程
	lock     mutex          // 互斥锁
}


go/src/runtime/chan.go
func makechan(t *chantype, size int64) *hchan {
make(chan struct{}, 10)的实现


go/src/runtime/chan.go
发送数据chan
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {


接收数据
go/src/runtime/chan.go
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {


chan被关闭
func closechan(c *hchan) {//close空指针的chan会崩溃  close已经close的chan也会崩溃