package main

import (
	"bufio"
	"fmt"
	"os"
)

// go run hello.go
func main() {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Please input your name:")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Found an error: %s\n", err)
	} else {
		// 去掉最后一个字节\n
		input = input[:len(input)-1]
		fmt.Printf("Hello, %s!\n", input)
	}
}
