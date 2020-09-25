package base

import (
	"fmt"
	"testing"
	"time"
)

func TestDurationBase(t *testing.T) {
	fmt.Println("sleep begin ", 1)
	// 正确，应该是默认转换
	//time.Sleep(1 * time.Second)

	// mismatched types int and Duration，类型乘一致
	//var num = int(1)

	// 可以
	var num = time.Duration(1)
	time.Sleep(num * time.Second)

	fmt.Println("sleep begin ", 2)
}
