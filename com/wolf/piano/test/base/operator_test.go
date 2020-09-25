package base

import (
	"fmt"
	"testing"
)

// 取反/异或
func TestUsingXor(t *testing.T) {
	var d uint8 = 2
	fmt.Printf("%08b\n", d)  // 00000010
	fmt.Printf("%08b\n", ^d) // 11111101
}

// 异或：不同则1
func TestUsingNot(t *testing.T) {
	var a uint8 = 0x82
	var b uint8 = 0x02
	fmt.Printf("%08b [A]\n", a)
	fmt.Printf("%08b [B]\n", b)

	fmt.Printf("%08b (NOT B)\n", ^b)
	fmt.Printf("%08b ^ %08b = %08b [B XOR 0xff]\n", b, 0xff, b^0xff)

	fmt.Printf("%08b ^ %08b = %08b [A XOR B]\n", a, b, a^b)
	fmt.Printf("%08b & %08b = %08b [A AND B]\n", a, b, a&b)
	fmt.Printf("%08b &^%08b = %08b [A 'AND NOT' B]\n", a, b, a&^b) // 不同位才取1
	fmt.Printf("%08b&(^%08b)= %08b [A AND (NOT B)]\n", a, b, a&(^b))
}
