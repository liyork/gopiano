package main

import "sync"

type Socket struct {
}

func (socket Socket) Send(bytes []byte) {

}

type SocketMutex struct {
	sync.Mutex
	socket *Socket
}

type SocketMap struct {
	sockets sync.Map
}

// 不想多次创建，但是socket初始化耗时，不能直接给返回了
// 写入新值的操作不光是调用一个api创建socket就完了，还要有一系列的初始化操作，我们必须保证在初始化完成之前，其他通过Load拿到这个实例的协程无法真正访问socket实例。
func (pushList *SocketMap) push(ip string, data []byte) {
	type SocketFunc func() *SocketMutex
	// 方法每次进来都是新的
	var w sync.WaitGroup
	// 闭包
	socketMutex := &SocketMutex{}
	w.Add(1)
	// 包装一层函数，让取到值，但是并没有初始化好的socker的协程进行等待
	socketInter, ok := pushList.sockets.LoadOrStore(ip, func() *SocketMutex {
		w.Wait()
		return socketMutex
	})
	if !ok { // store succ
		socketMutex.socket = NewSocket()
		//do some initial operation like connect
		// 初始化完则进行done
		w.Done()
	} else {
		i := socketInter.(SocketFunc)
		// 调用函数阻塞直到done，初始化完成
		socketMutex = i()
	}

	socketMutex.Lock()
	defer socketMutex.Unlock()
	// 已初始化好的
	socketMutex.socket.Send(data)
}

func NewSocket() *Socket {
	return nil
}
