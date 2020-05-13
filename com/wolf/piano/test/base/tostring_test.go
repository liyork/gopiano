package base

import (
	"fmt"
	"strconv"
	"testing"
)

type Person struct {
	name string
	sex  string
	age  int
}

type Person2 struct {
	name string
	sex  string
	age  int
}

// toString
func (p *Person) String() string {
	return p.name + "_" + strconv.Itoa(p.age)
}

func TestToString(t *testing.T) {
	person := &Person{"abc", "man", 11}
	fmt.Println(person)
}
