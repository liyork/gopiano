package base

import "strconv"

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
func (p Person) String() string {
	return p.name + "_" + strconv.Itoa(p.age)
}
