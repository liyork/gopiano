package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

// 在TCP/IP协议规定了在网络上必须采用网络字节顺序，也就是大端模式

func TestByteBase1(t *testing.T) {
	// 二进制形式：0000 0000 0000 0000 0001 0002 0003 0004
	var testInt int32 = 0x01020304 // 十六进制表示
	fmt.Printf("%d use big endian: \n", testInt)

	var testBytes []byte = make([]byte, 4)
	binary.BigEndian.PutUint32(testBytes, uint32(testInt)) //大端序模式
	fmt.Println("int32 to bytes:", testBytes)

	convInt := binary.BigEndian.Uint32(testBytes) //大端序模式的字节转为int32
	fmt.Printf("bytes to int32: %d\n\n", convInt)
}

func TestByteBase2(t *testing.T) {
	var v2 uint32
	var b2 [4]byte

	//内存地址 x00     x01          x02        x03
	// | 00000000 | 00000000 | 00000001 | 00000001 |
	// | b2[0]    | b2[1]    | b2[2]    | b2[3]    |
	v2 = 257 // 0x101,从左到右称为高位字节序->低位字节序
	fmt.Printf("v2:%b, %x\n", v2, v2)

	b2[3] = uint8(v2) // 直接转成uint8后等于 1，丢掉前24位
	fmt.Println("b2:", b2)
	b2[2] = uint8(v2 >> 8) // 移动后强转仅留8bit
	fmt.Println("move b2 8:", b2)
	b2[1] = uint8(v2 >> 16) // 移动16bit后
	fmt.Println("move b2 16:", b2)
	b2[0] = uint8(v2 >> 24)
	fmt.Println("move b2 32:", b2)
}

func TestBigEndianInternal(t *testing.T) {
	buf := make([]byte, 100)
	// 0 = {uint8} 0
	// 1 = {uint8} 0
	// 2 = {uint8} 1
	// 3 = {uint8} 47

	//  内部,大端模式
	// 	b[0] = byte(v >> 24) --[24~32)，v从左到右第一个8位
	//	b[1] = byte(v >> 16) --[16~24)，v从左到右第二个8位
	//	b[2] = byte(v >> 8)  --[8~16)，v从左到右第三个8位
	//	b[3] = byte(v)       --[0~8)，v从左到右第四个8位
	binary.BigEndian.PutUint32(buf[0:], uint32(303)) // 正常303对应二进制:100101111
}

//大端--整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//大端--字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
