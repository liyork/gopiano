package base

import (
	"fmt"
	"testing"
)

const baseIndex = 1 //起始索引

// 自增量——iota从0开始，每行+1(不包括空行和注释行,"_"不是空行)
// 在一个用 () 包围的常量定义列表中，所有常量的表达式可以省略，除了第一个
// 省略了表达式的常量使用其上一行常量的表达式，一直向上直到有值
type State int

const (
	Running State = iota + baseIndex
	Stopped
	Rebooting
	Terminated
)

func (s State) String() string {
	switch s {
	case Running:
		return "Running"
	case Stopped:
		return "Stopped"
	case Rebooting:
		return "Rebooting"
	case Terminated:
		return "Terminated"
	default:
		return "Unknown"
	}
}

func TestEnumBase(t *testing.T) {
	state := Running

	// print: "state 0"
	fmt.Println("state ", state)

	type T struct {
		Name  string
		Port  int
		State State
	}

	t1 := T{Name: "example", Port: 6666}

	// 输出： "t {Name:example Port:6666 State:Unknown}"
	fmt.Printf("t %+v\n", t1)
}
