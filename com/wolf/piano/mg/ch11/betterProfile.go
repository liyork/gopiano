package main

import (
	"fmt"
	"github.com/pkg/profile"
)

var VARIABLE int

func N11(n int) bool {
	for i := 2; i < n; i++ {
		if (n % i) == 0 {
			return false
		}
	}
	return true
}

func Multiply(a, b int) int {
	if a == 1 {
		return b
	}
	if a == 0 || b == 0 {
		return 0
	}
	if a < 0 {
		return -Multiply(-a, b)
	}
	return b + Multiply(a-1, b)
}

// go get github.com/pkg/profile
// go run betterProfile.go
func main() {
	defer profile.Start(profile.ProfilePath("./")).Stop()
	// 这个由 Dave Cheney 开发的 `github.com/pkg/profile` 包需要您插入一行声明来开启 **CPU 分析** 在您的 Go 应用中。如果您想开启 **内存分析**的话，您应该插入如下声明来代替：
	// defer profile.Start(profile.MemProfile).Stop()

	total := 0
	for i := 2; i < 200000; i++ {
		n := N11(i)
		if n {
			total++
		}
		fmt.Println("Total: ", total)
	}
	total = 0
	for i := 0; i < 5000; i++ {
		for j := 0; j < 400; j++ {
			k := Multiply(i, j)
			VARIABLE = k
			total++
		}
	}
	fmt.Println("Total: ", total)
}
