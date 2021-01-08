package v1

import (
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

// 用原子操作替换
// 原子操作适用于这里：原子读写数据，cas进行增加减少(因为最终是要增加减少的)

// 表示数据文件的接口类型
type DataFile interface {
	// 读取一个数据块
	Read() (rsn int64, d Data, err error)
	// 写入一个数据块
	Write(d Data) (wsn int64, err error)
	// 获取最后读取的数据块的序列号,Reading Serial Number，第几个块
	RSN() int64
	// 获取最后写入的数据块的序列号,Writing Serial Number
	WSN() int64
	// 获取数据块的长度
	DataLen() uint32
	// 关闭数据文件
	Close() error
}

// 表示数据的类型
type Data []byte

// 表示数据文件的实现类型
type myDataFile struct {
	roffset int64 // 读操作用到的偏移量

	woffset int64 // 写操作用到的偏移量

	fmutex sync.RWMutex // 用于文件的读写锁
	rcond  *sync.Cond
	f      *os.File // 文件

	dataLen uint32 // 数据块长度
}

// 新建一个数据文件的实例
func NewDataFile(path string, dataLen uint32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if dataLen == 0 {
		return nil, errors.New("Invalid data length!")
	}
	df := &myDataFile{f: f, dataLen: dataLen}
	df.rcond = sync.NewCond(df.fmutex.RLocker())
	return df, nil
}

// rsn:本次读取的第几块，d:读取的数据
func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	// 读取并更新读偏移量
	var offset int64
	for {
		offset := atomic.LoadInt64(&df.roffset)
		if atomic.CompareAndSwapInt64(&df.roffset, offset, df.roffset+int64(df.dataLen)) {
			break
		}
	}

	// 读取一个数据块
	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)
	df.fmutex.RLock()
	defer df.fmutex.RUnlock()
	for {
		// todo 这里并发读没事?也许是不同offset?
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err == io.EOF {
				df.rcond.Wait()
				continue
			}
			return // 不为EOF的错误，放弃再次尝试，立即返回给调用方，交给上层处理。
		}
		d = bytes
		return
	}
}

func (df *myDataFile) Write(d Data) (wsn int64, err error) {
	// 读取并更新写偏移量
	var offset int64
	for {
		offset = atomic.LoadInt64(&df.woffset)
		if atomic.CompareAndSwapInt64(&df.woffset, offset, df.woffset+int64(df.dataLen)) {
			break
		}
	}

	// 写入一个数据块
	wsn = offset / int64(df.dataLen)
	var bytes []byte
	if len(d) > int(df.dataLen) { // 截断
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}

	// Signal和Broadcast调用之前，无需锁定与之关联的锁。
	// 下面的锁定和解锁与Signal无联系
	df.fmutex.Lock()
	defer df.fmutex.Unlock()
	_, err = df.f.Write(bytes)
	// todo 而且这里仅仅通知一个可行吗?读是并发的，可能同时有很多等待，不过这里写入一次也只能写入一个块，通知多了也没用
	df.rcond.Signal()
	return
}

func (df *myDataFile) RSN() int64 {
	return atomic.LoadInt64(&df.roffset) / int64(df.dataLen)
}

func (df *myDataFile) WSN() int64 {
	return atomic.LoadInt64(&df.woffset) / int64(df.dataLen)
}

func (df *myDataFile) DataLen() uint32 {
	return df.dataLen
}

func (df *myDataFile) Close() error {
	return df.f.Close()
}
