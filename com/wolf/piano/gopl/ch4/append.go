package main

import "fmt"

func appendInt(x []int, y int) []int {
	var z []int
	// 增长一个后的长度
	zlen := len(x) + 1
	if zlen <= cap(x) {
		z = x[:zlen]
	} else { // 扩容
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x)
	}
	// 放入新值y
	z[len(x)] = y
	return z
}

func main() {
	var x, y []int
	for i := 0; i < 10; i++ {
		y = appendInt(x, i)
		fmt.Printf("%d cap=%d\t%v\n", i, cap(y), y)
		x = y
	}
}
