package other

import (
	"fmt"
	"testing"
)

func TestFib(t *testing.T) {
	fmt.Println(fib(3))
}

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
