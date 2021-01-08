package main

import "fmt"

// 使用copy()时你应小心翼翼，因为内建函数copy(dst,src)会以len(dst)和len(src)中的最小值为复制长度

func main() {
	a6 := []int{-10, 1, 2, 3, 4, 5}
	a4 := []int{-1, -2, -3, -4}
	fmt.Println("a6:", a6)
	fmt.Printf("a4:%v\n", a4)

	// dst,src
	copy(a6, a4)
	fmt.Println("a6:", a6)
	fmt.Printf("a4:%v\n", a4)
	fmt.Println()

	b6 := []int{-10, 1, 2, 3, 4, 5}
	b4 := []int{-1, -2, -3, -4}
	fmt.Println("b6:", b6)
	fmt.Printf("b4:%v\n", b4)

	// dst,src
	copy(b4, b6)
	fmt.Println("b6:", b6)
	fmt.Printf("b4:%v\n", b4)
	fmt.Println()

	fmt.Println()
	array4 := [4]int{4, -4, 4, -4}
	s6 := []int{1, -1, 1, -1, -5, 5}

	copy(s6, array4[0:])
	fmt.Println("array4:", array4[0:])
	fmt.Printf("s6:%v\n", s6)
	fmt.Println()

	array5 := [5]int{5, -5, 5, -5, 5}
	s7 := []int{7, 7, -7, 7, -7, 7}
	copy(array5[0:], s7)
	fmt.Println("array5:", array5)
	fmt.Printf("s7:%v\n", s7)
	fmt.Println()
}
