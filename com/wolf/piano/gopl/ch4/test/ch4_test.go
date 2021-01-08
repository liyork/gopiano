package test

import (
	"fmt"
	"testing"
)

func TestBase(t *testing.T) {
	var a [3]int
	fmt.Println(a[0])
	fmt.Println(a[len(a)-1])
	for i, v := range a {
		fmt.Printf("%d %d\n", i, v)
	}

	var q [3]int = [3]int{1, 2, 3}
	var r [3]int = [3]int{1, 2}
	fmt.Println(q, r[2])
	q2 := [...]int{1, 2, 3}
	fmt.Printf("%T\n", q2)

	r2 := [...]int{99: -1}
	fmt.Println(len(r2))

	month := [...]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
	Q2 := month[4:7]
	summer := month[6:9]
	fmt.Println(Q2)
	fmt.Println(summer)
	for _, s := range summer {
		for _, q := range Q2 {
			if s == q {
				fmt.Printf("%s appears in both\n", s)
			}
		}
	}

	var runes []rune
	for _, r := range "Hello, 世界" {
		runes = append(runes, r)
	}
	fmt.Printf("%q\n", runes)

}
