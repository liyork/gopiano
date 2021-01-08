package main

import "fmt"

var pc [256]byte

func init() {
	for i := range pc {
		//赋值给第i个元素
		// 第i下标除2，元素取值 + 1
		pc[i] = pc[i/2] + byte(1&1)
	}
}

func PopCount(x uint64) int {
	// 取x的最后8位
	return int(pc[byte(x>>0*8)] +
		//取次最后8位
		pc[byte(x>>1*8)] +
		pc[byte(x>>2*8)] +
		pc[byte(x>>3*8)] +
		pc[byte(x>>4*8)] +
		pc[byte(x>>5*8)] +
		pc[byte(x>>6*8)] +
		pc[byte(x>>7*8)])
}

func main() {
	fmt.Println("byte(1&1):", byte(1&1))
	fmt.Println("len:", len(pc))
	index := 0
	for i := range pc {
		if index%8 == 0 {
			fmt.Println()
		}
		fmt.Printf("%v ", pc[i])
		index++
	}

	//PopCount(1)
}
