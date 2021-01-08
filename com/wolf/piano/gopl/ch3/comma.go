package main

import "fmt"

// 三位一逗号，从后向前添加逗号，每次迭代下次都减少后三位
func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

func main() {
	fmt.Println(comma("1234567"))
}
