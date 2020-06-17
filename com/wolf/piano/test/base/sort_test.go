package base

import (
	"fmt"
	"sort"
	"testing"
)

type StuScore struct {
	name  string
	score int
}

type StuScores []StuScore

func (s StuScores) Len() int {
	return len(s)
}

func (s StuScores) Less(i, j int) bool {
	// 升序
	return s[i].score < s[j].score
}

func (s StuScores) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func TestSortStruct(t *testing.T) {
	stus := StuScores{
		{"alan", 95},
		{"hikerell", 91},
		{"acmfly", 96},
		{"leao", 90}}

	for _, v := range stus {
		fmt.Printf("k:%s,v:%d ", v.name, v.score)
	}
	fmt.Println()

	sort.Sort(stus)
	for _, v := range stus {
		fmt.Printf("k:%s,v:%d ", v.name, v.score)
	}

	fmt.Println("IS Sorted?", sort.IsSorted(stus))

	sort.Sort(sort.Reverse(stus))
	for _, v := range stus {
		fmt.Printf("k:%s,v:%d ", v.name, v.score)
	}
}

//sort包原生支持[]int、[]float64和[]string三种内建数据类型切片的排序操作
func TestInnerType(t *testing.T) {
	s := []int{5, 2, 6, 3, 1, 4}
	sort.Ints(s)
	fmt.Println(s)
	fmt.Println(sort.SearchInts(s, 3))

	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	fmt.Println(s)

}

func TestSearch(t *testing.T) {
	x := 11
	s := []int{3, 6, 8, 11, 45} //升序
	pos := sort.Search(len(s), func(i int) bool {
		fmt.Printf("i:%d, s[i]:%d\n", i, s[i])
		return s[i] >= x
	})
	if pos < len(s) && s[pos] == x {
		fmt.Println(x, "在s中的位置为：", pos)
	} else {
		fmt.Println("s不包含元素", x)
	}
}
