package concurrent

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// ======无缓冲通道的特点是，发送的数据需要被读取后，发送才会完成，它阻塞场景：

// 阻塞场景1:通道中无数据，由于没有写入协成，读会阻塞
func TestReadNoDataFromNoBufCh(t *testing.T) {
	noBufCh := make(chan int)
	<-noBufCh
	fmt.Println("read from no buffer channel success")
}

// 阻塞场景2:通道中无数据，由于无协程读取，向通道写会阻塞
func TestWriteNoBufCh(t *testing.T) {
	ch := make(chan int)
	ch <- 1
	fmt.Println("write success no block")
}

// ======有缓存通道的特点是，有缓存时可以向通道中写入数据后直接返回，缓存中有数据时可以从通道中读到数据直接返回，这时有缓存通道是不会阻塞的，
// 关键还是缓存中预留了数据

// 阻塞场景1:通道的缓存无数据，但执行读通道。
func ReadNoDataFromBufCh() {
	bufCh := make(chan int, 1)
	<-bufCh
	fmt.Println("read from no buffer channel success")
}

// 阻塞场景2:通道的缓存已经占满，向通道写数据，但无协程读。
func WriteBufChButFull() {
	ch := make(chan int, 1)
	// make ch full
	ch <- 100
	ch <- 1
	fmt.Println("write success no block")
}

// ===== 使用Select实现无阻塞读写
// 原理:select是执行选择操作的一个结构，它里面有一组case语句，它会执行其中无阻塞的那一个，如果都阻塞了，那就等待其中一个不阻塞，进而继续执行，
// 它有一个default语句，该语句是永远不会阻塞的，借助它实现无阻塞的操作。在通道不可读或者不可写的时候，不再阻塞等待，而是直接返回。

func ReadWithSelect(ch chan int) (x int, err error) {
	select {
	case x = <-ch:
		return x, nil
	default:
		return 0, errors.New("channel has no data")
	}
}

// 无缓冲通道,读
func TestReadNoDataFromNoBufChWithSelect(t *testing.T) {
	bufCh := make(chan int)
	if v, err := ReadWithSelect(bufCh); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("read: %d\n", v)
	}
}

// 有缓冲通道,读
func TestReadNoDataFromBufChWithSelect(t *testing.T) {
	bufCh := make(chan int, 1)
	if v, err := ReadWithSelect(bufCh); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("read: %d\n", v)
	}
}

func WriteChWithSelect(ch chan int) error {
	select {
	case ch <- 1:
		return nil
	default:
		return errors.New("channel blocked, can not write")
	}
}

// 无缓冲通道,写
func TestWriteNoBufChWithSelect(t *testing.T) {
	ch := make(chan int)
	if err := WriteChWithSelect(ch); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("write success")
	}
}

// 有缓冲通道,写
func TestWriteBufChButFullWithSelect(t *testing.T) {
	ch := make(chan int, 1)
	// make ch full
	ch <- 100
	if err := WriteChWithSelect(ch); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("write success")
	}
}

// ====== 使用Select+超时改善无阻塞读写
// 使用default实现的无阻塞通道阻塞有一个缺陷：当通道不可读或写的时候，会立即返回。
// 实际场景中更多的需求是，希望尝试读一会数据，或者尝试写一会数据，如果实在没法读写，再返回，程序继续做其它的事情。
// 使用定时器替代default可以解决这个问题。比如，我给通道读写数据的容忍时间是500ms，如果依然无法读写，就即刻返回

func ReadWithSelectTimeout(ch chan int) (x int, err error) {
	timeout := time.NewTimer(time.Microsecond * 500)

	select {
	case x = <-ch:
		return x, nil
	case <-timeout.C:
		return 0, errors.New("read time out")
	}
}

// 无缓冲通道,读超时
func TestReadNoDataFromNoBufChWithSelectTimeout(t *testing.T) {
	bufCh := make(chan int)
	if v, err := ReadWithSelectTimeout(bufCh); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("read: %d\n", v)
	}
}

func WriteChWithSelectTimeout(ch chan int) error {
	timeout := time.NewTimer(time.Microsecond * 500)

	select {
	case ch <- 1:
		return nil
	case <-timeout.C:
		return errors.New("write time out")
	}
}

// 无缓冲通道,写超时
func TestWriteNoBufChWithSelectTimeout(t *testing.T) {
	ch := make(chan int)
	if err := WriteChWithSelectTimeout(ch); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("write success")
	}
}
