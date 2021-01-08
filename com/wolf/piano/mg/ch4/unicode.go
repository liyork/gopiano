package main

import (
	"fmt"
	"unicode"
)

// go run unicode.go
func main() {
	const sL = "\x99\x00ab\x50\x00\x23\x50\x29\x9c"

	for i := 0; i < len(sL); i++ {
		// 检查字符串`sL`的每个元素是否是rune类型，如果是的话将返回`true`否则返回false
		if unicode.IsPrint(rune(sL[i])) {
			fmt.Printf("%c\n", sL[i])
		} else {
			fmt.Println("Not printable!")
		}
	}
}
