package buffer

import (
	syserrors "errors"
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/gopcp/ch6/webcrawler/errors"
	"sync"
	"sync/atomic"
)

// FIFO的缓冲器的接口类型，包装chan，解决ch的相关panic问题
type Buffer interface {
	// 获取本缓冲器的容量
	Cap() uint32
	// 获取本缓冲器中数据的数量
	Len() uint32
	// 向缓冲器放入数据，非阻塞，若缓冲器已关闭，则会直接返回非nil的错误值
	Put(datum interface{}) (bool, error)
	// 从缓冲器中获取，非阻塞，若缓冲器已关闭，则会直接返回非nil的错误值
	Get() (interface{}, error)
	// 关闭缓冲器，若之前已关闭则false，否则true
	Close() bool
	// 判断缓冲器是否已关闭
	Closed() bool
}

// 缓冲器接口的实现类型
type myBuffer struct {
	// 存放数据的通道
	ch chan interface{}
	// 关闭状态：0-未关闭，1-已关闭
	closed uint32
	// 为了消除因关闭缓冲器而产生的竞态条件的读写锁
	closingLock sync.RWMutex
}

// size代表缓冲器的容量
func NewBuffer(size uint32) (Buffer, error) {
	if size == 0 {
		errMsg := fmt.Sprintf("illegal size for buffer: %d", size)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	return &myBuffer{
		ch: make(chan interface{}, size),
	}, nil
}

func (buf *myBuffer) Cap() uint32 {
	return uint32(cap(buf.ch))
}

func (buf *myBuffer) Len() uint32 {
	return uint32(len(buf.ch))
}

// 非阻塞
func (buf *myBuffer) Put(datum interface{}) (ok bool, err error) {
	buf.closingLock.RLock()
	defer buf.closingLock.RUnlock() // 保护下面全部，若只检查这里，执行下面时若close了那么写入closed的chan也会panic
	if buf.Closed() {               // 已关闭
		return false, ErrClosedBuffer
	}
	select {
	case buf.ch <- datum:
		ok = true
	default: // 满了
		ok = false
	}
	return
}

// 表示缓冲器已关闭的错误
var ErrClosedBuffer = syserrors.New("closed buffer")

// 非阻塞
func (buf *myBuffer) Get() (interface{}, error) {
	select {
	case datum, ok := <-buf.ch:
		if !ok {
			return nil, ErrClosedBuffer
		}
		return datum, nil
	default:
	}
	return nil, nil
}

func (buf *myBuffer) Close() bool {
	if atomic.CompareAndSwapUint32(&buf.closed, 0, 1) { // 控制并发调用Close，只有一个成功
		buf.closingLock.Lock() // 和Put保持互斥和可见性
		close(buf.ch)          // 针对此操作需要加锁
		buf.closingLock.Unlock()
		return true
	}
	return false
}

func (buf *myBuffer) Closed() bool {
	if atomic.LoadUint32(&buf.closed) == 0 {
		return false
	}
	return true
}
