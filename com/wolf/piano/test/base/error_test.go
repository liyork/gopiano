package base

import (
	"errors"
	"fmt"
	"testing"
)

//错误接口
//type error interface {
//  Error() string
//}

// 对于返回值传递是值传递
func TestErrorAddress(t *testing.T) {
	fmt.Println(11111)
	e := occurError()
	fmt.Println("e:", e)
	fmt.Println("&e2:", &e)
}

func occurError() error {
	fmt.Println("222")
	e := errors.New("error happen")
	fmt.Println("&e1", &e)
	return e
}

// 对于返回指针传递是指针传递
func TestErrorPointAddress(t *testing.T) {
	fmt.Println(11111)
	e := occurErrorPoint()
	fmt.Println("e:", e)
	fmt.Println("&e2:", e)
}

func occurErrorPoint() *error {
	fmt.Println("222")
	e := errors.New("error happen")
	fmt.Println("&e1", &e)
	return &e
}
