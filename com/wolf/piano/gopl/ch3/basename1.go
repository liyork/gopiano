package main

import "fmt"

// 取最后一个/之后的，然后排除.之后的内容
func basename(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}

func main() {
	println(basename("a/b/c.go"))
	println(basename("c.d.go"))
	println(basename("abc"))
	fmt.Println("")
	println(basename2("a/b/c.go"))
	println(basename2("c.d.go"))
	println(basename2("abc"))
}
