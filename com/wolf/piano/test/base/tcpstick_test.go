package base

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"testing"
)

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn) //读取conn里接收的内容
	var buf = make([]byte, 1024)
	for {
		n, err := reader.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		buf = make([]byte, 1024)
		fmt.Println("收到client发来的数据：", recvStr)
	}
}

func process2(conn net.Conn) {
	defer conn.Close()
	//reader := bufio.NewReader(conn)
	//reader := bufio.NewReaderSize(conn, 15) // buf的大小，若不满足单次消息则就会一直收不到
	reader := bufio.NewReaderSize(conn, 31) // buf的大小，若不满足单次消息则就会一直收不到
	for {
		msg, err := Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		fmt.Println("收到client发来的数据：", msg)
	}
}

func Encode(message string) ([]byte, error) {
	// 读取消息的长度，转换成int32类型--固定（占4个字节）
	var length = int32(len(message))
	var out = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(out, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(out, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

const headerLen = 4

func Decode(reader *bufio.Reader) (string, error) {
	// 读取消息的长度
	lengthByte, _ := reader.Peek(headerLen) // 读取前4个字节的数据

	// not enough to headerlen
	if int32(reader.Buffered()) < headerLen {
		return "", nil
	}

	lengthBuff := bytes.NewBuffer(lengthByte)
	var bodyLen int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &bodyLen)
	if err != nil {
		return "", err
	}
	fmt.Println("before read, Buffered len", reader.Buffered())
	// 若小于单次消息整体长度则返回，等待下次读好了再处理
	if int32(reader.Buffered()) < headerLen+bodyLen {
		return "", err
	}

	// 读取真正的单条消息数据
	pack := make([]byte, int(headerLen+bodyLen))
	readLen, err := reader.Read(pack)
	fmt.Println("after read, Buffered len", reader.Buffered(), " readlen:", readLen)
	if err != nil {
		return "", err
	}
	// 仅返回body
	return string(pack[headerLen:]), nil
}

// ======stick

func TestStickServer(t *testing.T) {
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)
	}
}

func TestStickClient(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		conn.Write([]byte(msg))
	}
}

// ======nonstick

func TestNonStickServer(t *testing.T) {
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process2(conn)
	}
}

func TestNonStickClient(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		data, err := Encode(msg)
		if err != nil {
			fmt.Println("encode msg failed, err:", err)
			return
		}
		conn.Write(data)
	}
}
