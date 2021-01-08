package test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestBase(t *testing.T) {
	got := strings.Split("a:b:c", ":")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		fmt.Println("not deep equal")
	}

	var a, b []string = nil, []string{}
	fmt.Println(reflect.DeepEqual(a, b))

	var c, d map[string]int = nil, make(map[string]int)
	fmt.Println(reflect.DeepEqual(c, d))
}
