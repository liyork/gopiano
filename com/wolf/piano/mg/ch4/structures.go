package main

import "fmt"

// go run structures.go
func main() {
	type XYZ struct {
		X int
		Y int
		Z int
	}

	var s1 XYZ
	fmt.Println(s1.Y, s1.Z)

	p1 := XYZ{23, 12, -2}
	p2 := XYZ{Z: 12, Y: 13}
	fmt.Println(p1)
	fmt.Println(p2)

	pSlice := [4]XYZ{}
	// 当将一个结构体分配给结构体数组，那么这个结构体就会被深拷贝至这个数组，意味着改变原结构体是对数组中的结构体没有影响的
	pSlice[2] = p1
	pSlice[0] = p2
	fmt.Println(pSlice)

	//p2 = XYZ{1, 2, 3}
	p2.X = 4
	fmt.Println(p2, " | ", pSlice)
}
