package main

import (
	"fmt"
	"net"
	"os"
)

//Go语言标准库对此过程进行了抽象和封装。无论我们期望使用什么协议建立什么形式的连接，都只需要调用net.Dial()即可
//Dial（）函数的
//net参数是网络协议的名字，addr参数是IP地址或域名，而端口号以“：”的形式跟随在地址或域名的后面，端口号可选。
//tcp、tcp4、tcp6、udp、udp4、udp6、ip、ip4、ip6、ip4:icmp、ip4:1

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "host")
		os.Exit(1)
	}
	service := os.Args[1]
	conn, err := net.Dial("ip4:icmp", service)
	checkError(err)

	var msg [512]byte
	msg[0] = 8  // echo
	msg[1] = 0  // code 0
	msg[2] = 0  // checksum
	msg[3] = 0  // checksum
	msg[4] = 0  // identifier[0]
	msg[5] = 13 // identifier[1]
	msg[6] = 0  // sequence[0]
	msg[7] = 37 // sequence[1]
	len := 8
	check := checkSum(msg[0:len])

	// 高8位
	msg[2] = byte(check >> 8)
	// 低8位
	msg[3] = byte(check & 255)
	_, err = conn.Write(msg[0:len])
	checkError(err)

	_, err = conn.Read(msg[0:])
	checkError(err)
	fmt.Println("Got response")
	if msg[5] == 13 {
		fmt.Println("Identifier matches")
	}
	if msg[7] == 37 {
		fmt.Println("Sequence matches")
	}
	os.Exit(0)

}
func checkSum(msg []byte) uint16 {
	sum := 0
	// 先假设为偶数
	for n := 1; n < len(msg)-1; n += 2 {
		sum += int(msg[n])*256 + int(msg[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
