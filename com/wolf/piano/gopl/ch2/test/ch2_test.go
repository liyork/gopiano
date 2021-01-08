package test

import (
	"fmt"
	"testing"
)

func TestBase(t *testing.T) {
	var s string
	fmt.Println(s)

	i, j := 0, 1
	i, j = j, i
	fmt.Println(i, j)

	x := 1
	p := &x
	fmt.Println(*p)
	*p = 2
	fmt.Println(x)

	var x1, y int
	fmt.Println(&x1 == &x1, &x1 == &y, &x1 == nil)

	fmt.Println(f() == f())

	p1 := new(int)
	fmt.Println(*p1)
	*p1 = 2
	fmt.Println(*p1)

}

func f() *int {
	v := 1
	return &v
}
