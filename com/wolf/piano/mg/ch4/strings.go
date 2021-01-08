package main

import "fmt"

// go run strings.go
func main() {
	//每一个`\xAB`都代表一个字符
	const sLiteral = "\x99\x42\x32\x55\x50\x35\x23\x50\x29\x9c"
	fmt.Println(sLiteral)
	// `%x`会得到`\xAB`中的`AB`部分
	fmt.Printf("x: %x\n", sLiteral)
	// 得到`sLiteral`的字符数量
	fmt.Printf("sLiteral length: %d\n", len(sLiteral))

	// 你可以像操作一个切片那样去操作字符串。
	// %x输出字符，如果你想打印字符串格式，就要使用`%s`
	for i := 0; i < len(sLiteral); i++ {
		fmt.Printf("%x ", sLiteral[i])
	}
	fmt.Println()

	// 使用`%q`作为字符串格式化参数，可以安全地打印出带双引号的字符串，
	fmt.Printf("q: %q\n", sLiteral)
	// `%+q`可以保证输出是ASCII格式
	fmt.Printf("+q: %+q\n", sLiteral)
	fmt.Printf(" x: % x\n", sLiteral)
	// 输出字符串
	fmt.Printf("s: As a string: %s\n", sLiteral)

	// 在这里定义了字符串`s2`，内容是3个unicode字符。使用`%#U`可以打印出`U+0058`格式的字符，
	s2 := "€£³"
	for x, y := range s2 { // `range`关键字能够迭代包含Unicode字符的字符串，这样就可以逐个处理Unicode字符。
		fmt.Printf("%#U starts at byte position %d\n", y, x)
	}

	// `s2`包含的是Unicode字符，Unicode字符的字节数量是大于该字符串中的元素数量的，而`len()`函数计算的是字节数量，所以`len(s2)`的值是7而不是3。
	fmt.Printf("s2 length: %d\n", len(s2))

	const s3 = "ab12AB"
	fmt.Println("s3:", s3)
	fmt.Printf("x: % x\n", s3)

	fmt.Printf("s3 length: %d\n", len(s3))

	for i := 0; i < len(s3); i++ {
		fmt.Printf("%x ", s3[i])
	}
	fmt.Println()
}
