package main

import (
	"errors"
	"fmt"
)

// 函数
func divide(dividend int, divisor int) (result int, err error) {
	if divisor == 0 {
		err = errors.New("division by zero")
		return
	}
	result = dividend / divisor
	return
}

// 定义二元操作的函数类型
type binaryOperation func(operand1 int, operand2 int) (result int, err error)

//把二元操作的实现权留给了operate函数的使用者。
func operate(op1 int, op2 int, bop binaryOperation) (result int, err error) {
	if bop == nil {
		err = errors.New("invalid binary operation funciton")
		return
	}
	return bop(op1, op2)
}

// 方法
type myInt int

// 接收者
func (i myInt) add(another int) myInt {
	// 产生新变量
	i = i + myInt(another)
	return i
}

// add方法是myInt的指针方法
func (i *myInt) add1(another int) myInt {
	*i = *i + myInt(another)
	return *i
}

func main() {
	result, err := operate(1, 2, func(op1 int, op2 int) (int, error) {
		return op1 + op2, nil
	})
	fmt.Println(result, err)

	result, err = operate(1, 2, divide)
	fmt.Println(result, err)

	i1 := myInt(1)
	i2 := i1.add(2)
	println(i1, i2)

	i2 = i1.add1(2)
	println(i1, i2)
}
