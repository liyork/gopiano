package base

import (
	"fmt"
	"testing"
)

func TestSelectBase(t *testing.T) {
	var a = ' '
	switch a {
	case ' ': // 默认带有break
		//fmt.Println("blank", )
	case '\t':
		fmt.Println("tttt")
	}
}

func TestSelectBase2(t *testing.T) {
	var a = ' '
	switch a {
	case ' ':
		fallthrough //强制执行下一个 case 代码块
	case '\t':
		fmt.Println("tttt")
	}
}
