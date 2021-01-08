package buffer

import (
	syserrors "errors"
	"sync"
	"sync/atomic"
)

// 数据缓冲池的接口类型
type Pool interface {
	// 获取缓冲池中缓冲器的统一容量(buffer的ch容量)
	BufferCap() uint32
	// 获取池中缓冲器的数量
	BufferNumber() uint32
	// 获取池中缓冲器的最大数量
	MaxBufferNumber() uint32
	// 获取缓冲池中数据的总数
	Total() uint64
	// 向缓冲池中放入数据，阻塞，若缓冲池已关闭，则直接返回非nil的错误值
	Put(datum interface{}) error
	// 从缓冲池获取数据，阻塞，若缓冲池已关闭，则直接返回非nil的错误值
	Get() (datum interface{}, err error)
	// 关闭缓冲池，若缓冲池之前已关闭则false，否则true
	Close() bool
	// 获取缓冲池是否已关闭
	Closed() bool
}

// 数据缓冲池接口的实现类型，通过chan进行存储buffer，先通过chan交互获取、放入buffer，然后通过buffer放入、获取数据
type myPool struct {
	// 缓冲器的统一容量
	bufferCap uint32
	// 缓冲器的实际数量
	bufferNumber uint32
	// 缓冲器的最大数量
	maxBufferNumber uint32
	// 池中数据的总数
	total uint64
	// 存放缓冲器的通道
	bufCh chan Buffer
	// 缓冲池的关闭状态：0-未关闭；1-已关闭
	closed uint32
	// 保护内部共享资源的读写锁
	rwlock sync.RWMutex
}

func (pool *myPool) BufferCap() uint32 {
	panic("implement me")
}

func (pool *myPool) Total() uint64 {
	panic("implement me")
}

func NewPool() Pool {
	return &myPool{}
}

var ErrClosedBufferPool = syserrors.New("ErrClosedBufferPool")

func (pool *myPool) Put(datum interface{}) (err error) {
	if pool.Closed() {
		return ErrClosedBufferPool
	}
	var count uint32
	maxCount := pool.bufferNumber * 5
	var ok bool
	// 不断从ch中获取buffer并试图放入，直到ok或失败超过次数停止
	for buf := range pool.bufCh {
		// 非阻塞
		ok, err = pool.putData(buf, datum, &count, maxCount)
		if ok || err != nil {
			break
		}
	}
	return
}

// 向给定的缓冲器放入数据，放入后，在必要时把缓冲器归还给池
func (pool *myPool) putData(buf Buffer, datum interface{}, count *uint32, maxCount uint32) (ok bool, err error) {
	if pool.Closed() {
		return false, ErrClosedBufferPool
	}
	// 用到了rwlock的读锁，因为这其中包含了向bufCh发送值的操作。若在方法即将执行结束时，发现缓冲池已关闭，那就不会归还拿到的缓冲器，同时把对应的错误赋值给err。
	// 这里用rwlock的读锁,应该是针对close进行保护，针对多读，大家个各自拿到不同buf，最后放入时只要保证没close，ch本身是并发安全的
	defer func() {
		pool.rwlock.RLock()
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			err = ErrClosedBufferPool
		} else {
			pool.bufCh <- buf
		}
		pool.rwlock.RUnlock()
	}()

	ok, err = buf.Put(datum)
	if ok {
		atomic.AddUint64(&pool.total, 1)
		return
	}
	// 缓冲器已关闭
	if err != nil {
		return
	}
	// 若因缓冲器已满而未放入数据，递增计数
	(*count)++

	// 若尝试向缓冲器放入数据的失败次数达到阈值，且池中缓冲器的数量未达到最大值，尝试创建一个新的缓冲器，先放入数据再把它放入池
	if *count >= maxCount && pool.BufferNumber() < pool.MaxBufferNumber() {
		pool.rwlock.Lock()                                // 发现不满足时，仅有一个协程准备扩容，其他的稍后进来再判断，用了写锁
		if pool.BufferNumber() < pool.MaxBufferNumber() { // 双检锁
			if pool.Closed() {
				pool.rwlock.Unlock()
				return
			}
			newBuf, _ := NewBuffer(pool.bufferCap)
			newBuf.Put(datum)
			pool.bufCh <- newBuf
			atomic.AddUint32(&pool.bufferNumber, 1)
			atomic.AddUint64(&pool.total, 1)
			ok = true
		}
		pool.rwlock.Unlock()
		// 及时清零“已满”计数可以有效减少不必要的操作和资源消耗。
		*count = 0
	}
	return
}

func (pool *myPool) Get() (datum interface{}, err error) {
	if pool.Closed() {
		return nil, ErrClosedBufferPool
	}
	var count uint32
	maxCount := pool.BufferNumber() * 10
	for buf := range pool.bufCh {
		datum, err = pool.getData(buf, &count, maxCount)
		if datum != nil || err != nil {
			break
		}
	}
	return
}

// 从给定的缓冲器获取数据，并在必要时把缓冲器归还给池
func (pool *myPool) getData(buf Buffer, count *uint32, maxCount uint32) (datum interface{}, err error) {
	if pool.Closed() {
		return nil, ErrClosedBufferPool
	}
	defer func() {
		// 若尝试从缓冲器获取数据的失败次数达到阈值，同时当前缓冲器已空且池中缓冲器的数量大于1，就直接关掉当前缓冲器，并不归还给池
		if *count >= maxCount && buf.Len() == 0 && pool.BufferNumber() > 1 {
			buf.Close()
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			*count = 0
			return
		}

		pool.rwlock.RLock() // 锁定rwlock的读锁，以避免向已关闭的pool发送值。只有close/扩容时用写锁，也就是读写都可以用读锁并发执行
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			err = ErrClosedBufferPool
		} else {
			pool.bufCh <- buf
		}
		pool.rwlock.RUnlock()

		datum, err = buf.Get()
		if datum != nil {
			atomic.AddUint64(&pool.total, ^uint64(0))
			return
		}
		if err != nil {
			return
		}
		// 若因缓冲器已空未取出数据，递增计数
		(*count)++
		return
	}()
	return
}

func (pool *myPool) Close() bool {
	if !atomic.CompareAndSwapUint32(&pool.closed, 0, 1) { // 防止并发关闭
		return false
	}

	pool.rwlock.Lock()
	defer pool.rwlock.Unlock()
	close(pool.bufCh)
	for buf := range pool.bufCh {
		buf.Close()
	}
	return true
}

func (pool *myPool) Closed() bool {
	return atomic.LoadUint32(&pool.closed) == 1
}

func (pool *myPool) BufferNumber() uint32 {
	return pool.bufferNumber
}

func (pool *myPool) MaxBufferNumber() uint32 {
	return pool.maxBufferNumber
}
