package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// `/dev/random`系统设备的目的是生成随机数据，你可以使用这些数据来测试程序，
// 在本例中，你将为随机数生成器生成随机数种子。
// ls -l /dev/random
// `/dev/random`返回二进制数据，需要解码

// go run devRandom.go
func main() {
	f, err := os.Open("/dev/random")
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	var seed int64
	binary.Read(f, binary.LittleEndian, &seed)
	fmt.Println("Seed:", seed)
}
