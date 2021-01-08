package main

import "fmt"

// go run runes.go
func main() {
	const r1 = '€'
	fmt.Println("(int32) r1:", r1)
	fmt.Printf("(HEX) r1: %x\n", r1)
	fmt.Printf("(as a String) r1: %s\n", r1)    // 字符串
	fmt.Printf("(as a character) r1: %c\n", r1) // 字符

	// 字节切片就是一系列runes的集合
	fmt.Println("A string is a collection of runes:", []byte("Mihalis"))
	aString := []byte("Mihalis")
	for x, y := range aString {
		fmt.Println(x, y)
		// `%c`可以将runes转换为字符输出；如果想要以字符串的形式输出字节数组，应使用`fmt.Printf()`结合`%s`
		fmt.Printf("Char: %c\n", aString[x])
		fmt.Printf("Str: %s\n", aString[x])
		fmt.Println("")
	}
	fmt.Printf("%s\n", aString)
}
