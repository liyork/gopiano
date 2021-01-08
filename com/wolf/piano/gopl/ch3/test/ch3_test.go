package test

import (
	"fmt"
	"strconv"
	"testing"
	"unicode/utf8"
)

func TestBase(t *testing.T) {
	var x uint8 = 1<<1 | 1<<5
	var y uint8 = 1<<1 | 1<<2

	fmt.Printf("%08b\n", x)
	fmt.Printf("%08b\n", y)
	fmt.Println("")
	fmt.Printf("%08b\n", x&y)
	fmt.Printf("%08b\n", x|y)
	fmt.Printf("%08b\n", x^y)
	// 将x中与y相异的部分保留，其他清零
	fmt.Printf("%08b\n", x&^y)

	for i := uint(0); i < 8; i++ {
		if x&(1<<i) != 0 {
			fmt.Println(i)
		}
	}
	fmt.Printf("%08b\n", x<<1)
	fmt.Printf("%08b\n", x>>1)

	o := 0666
	// 用第一个参数以8进制打印。  省略0
	fmt.Printf("%d %[1]o %#[1]o\n", o)
	z := int64(0xdeadbeef)
	// 用第一个参数以16进制打印。  省略0
	fmt.Printf("%d %[1]x %#[1]x %#[1]X\n", z)

	ascii := 'a'
	unicode := '国'
	newline := '\n'
	fmt.Printf("%d %[1]c %[1]q\n", ascii)
	fmt.Printf("%d %[1]c %[1]q\n", unicode)
	fmt.Printf("%d %[1]q\n", newline)

	s := "hello, world"
	fmt.Println(len(s))
	// 从0开始
	fmt.Println(s[0], s[7])
	fmt.Println(s[0:5])

	s = "Hello, 世界"
	fmt.Println(len(s))                    //7+6
	fmt.Println(utf8.RuneCountInString(s)) //每个rune

	for i, r := range "Hello, 世界" {
		fmt.Printf("%d %q %d\n", i, r, r)
	}

	s = "abc"
	b := []byte(s)
	s2 := string(b)
	fmt.Println(s, s2)

	x2 := 123
	y2 := fmt.Sprintf("%d", x)
	fmt.Println(y2, strconv.Itoa(x2))

	fmt.Println(strconv.FormatInt(int64(x2), 2))

	s = fmt.Sprintf("x=%b", x)
	fmt.Println(s)

	x3, err := strconv.Atoi("123")
	y3, err := strconv.ParseInt("123", 10, 64)
	fmt.Println(x3, y3, err)
}
