package test

import (
	"fmt"
	"testing"
)

func TestDebugBase1(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println("i1:", i)
	}
}

func TestCover1(t *testing.T) {
	test1("s1")
}
