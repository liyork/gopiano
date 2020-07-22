package base

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

// bufio 包实现了带缓存的 I/O 操作

func TestReaderSimple(t *testing.T) {
	data := "hello world"
	b := bufio.NewReader(strings.NewReader(data))
	if s := readBytes(b); s != "hello world" {
		t.Errorf("simple hello world test failed: got %q", s)
	}
}

func readBytes(buf *bufio.Reader) string {
	var b [1000]byte
	index := 0
	for {
		c, err := buf.ReadByte()
		if err == io.EOF {
			break
		}
		if err == nil {
			b[index] = c
			index++
		} else if err != iotest.ErrTimeout {
			panic("Data: " + err.Error())
		} else {
			fmt.Println("err occur when ReadByte, err:", err)
		}
	}
	return string(b[0:index])
}

func TestBuffered(t *testing.T) {
	data := "hello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello world" +
		"hello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello world"
	reader := strings.NewReader(data)
	newReader := bufio.NewReader(reader)
	fmt.Println("Buffered:", newReader.Buffered())
	// 读到内部b.buf(增长了b.w)然后返回b.buf[b.r : b.r+n]
	newReader.Peek(1)
	// b.w - b.r
	fmt.Println("Buffered:", newReader.Buffered())

	var a = make([]byte, 1)
	n, _ := newReader.Read(a)
	fmt.Println("read n:", string(a[0:n]))
	fmt.Println("Buffered:", newReader.Buffered())
}
