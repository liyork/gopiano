package base

import (
	"fmt"
	"testing"
)

func TestMapBase(t *testing.T) {
	maps := make(map[int]string, 3)
	maps[0] = "a"
	maps[1] = "b"
	maps[2] = "c"

	for k, v := range maps {
		fmt.Printf("k:%v,v:%v\n", k, v)
	}
}

type tmpStruct struct {
	name int
}

func TestMapPoint(t *testing.T) {
	maps := map[int]*tmpStruct{}

	s1 := maps[1]
	if s1 == nil {
		s1 = &tmpStruct{name: 1}
		maps[1] = s1
	}
	s1.name = 22222222

	// hashmap再扩容应该是不会对内容有影响只不过是重新hash定位而已
	for i := 1; i < 100; i++ {
		maps[i+1] = &tmpStruct{name: i + 1}
	}

	fmt.Println("maps[1]:", maps[1])
}

func TestDeleteNil(t *testing.T) {
	maps := map[string]map[int]*tmpStruct{}
	structs := maps["a"]
	delete(structs, 111)
}
