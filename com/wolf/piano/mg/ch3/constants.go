package main

import "fmt"

type Digit int
type Power2 int

const PI = 3.1415926

const (
	C1 = "C1C1C1"
	C2 = "C2C2C2"
	C3 = "C3C3C3"
)

// go run constants.go
func main() {
	const s1 = 123
	var v1 float32 = s1 * 12
	fmt.Println(v1)
	fmt.Println(PI)

	const (
		Zero Digit = iota // 从0开始
		One
		Two
		Three
		Four
	)
	// 类似于
	//const (
	//	Zero = 0
	//	One = 1
	//	Two = 2
	//	Three = 3
	//	Four = 4
	//)

	fmt.Println(One)
	fmt.Println(Two)

	const (
		p2_0 Power2 = 1 << iota
		_
		p2_2
		_ // _`符号是跳过本次常量声明并递增，iota的递增属性是可以用到表达式中的。
		p2_4
		_
		p2_6
	)

	fmt.Println("2^0:", p2_0)
	fmt.Println("2^2:", p2_2)
	fmt.Println("2^4:", p2_4)
	fmt.Println("2^6:", p2_6)

}
