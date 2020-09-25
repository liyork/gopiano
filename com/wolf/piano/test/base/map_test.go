package base

import (
	"fmt"
	"runtime/debug"
	"sort"
	"strings"
	"testing"
	"time"
)

// 无序
func TestMapBase(t *testing.T) {
	maps := make(map[int]string, 3)
	maps[0] = "a"
	maps[1] = "b"
	maps[2] = "c"

	for k, v := range maps {
		fmt.Printf("k:%v,v:%v\n", k, v)
	}
}

func TestMapRef(t *testing.T) {
	maps := make(map[int]string, 3)
	maps[0] = "a"
	maps[1] = "b"
	maps[2] = "c"

	// 两个引用一个对象
	map2 := maps
	map2[0] = "x"

	fmt.Println(maps)
}

func TestMapExist(t *testing.T) {
	maps := make(map[int]int, 3)
	maps[0] = 111

	// 不存在返回默认值
	fmt.Println("maps not exist return:", maps[1] == 0)
	// 区分是否存在exist
	if value, exist := maps[1]; !exist {
		fmt.Println("value not exist,value:", value)
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

func TestDelete(t *testing.T) {
	maps := map[int]string{1: "a", 2: "b"}
	delete(maps, 1)
	fmt.Println(maps)
}

// 对nil删除也安全
func TestDeleteNil(t *testing.T) {
	maps := map[string]map[int]*tmpStruct{}
	structs := maps["a"]
	delete(structs, 111)

	maps = nil
	delete(maps, "111")
}

func TestMapReAssign(t *testing.T) {
	maps := map[int]string{}
	maps[1] = "a"
	reAssignMap(&maps)
	fmt.Printf("maps1:%p\n", maps)
}

// 对map的指针地址内容进行重新设定
func reAssignMap(maps *map[int]string) {
	*maps = map[int]string{}
	fmt.Printf("maps2:%p\n", maps)
}

func TestMapPass(t *testing.T) {
	maps := map[int]string{}
	maps[1] = "a"
	passMap(maps)
	fmt.Printf("maps1:%p\n", maps)
	fmt.Printf("maps1:%v\n", maps)
}

// 传递的是map的地址值的拷贝，都是一样0xc00010c000，使用时都用这个地址操作那个地方的map
func passMap(maps map[int]string) {
	fmt.Printf("maps2:%p\n", maps)
	maps[1] = "b"
	fmt.Printf("maps2:%v\n", maps)
}

// GODEBUG=gctrace=1 go test -v -run TestMapLeak map_test.go
// 运行可见，即使代码逻辑清空了map，但进程内存使用并没有像预期那样“实报实销”：
// 解决：定期替换成新的map，释放旧的map对象。
func TestMapLeak(t *testing.T) {
	m := make(map[int]string)
	s := strings.Repeat("x", 1024)
	for i := 0; i < 100000000; i++ {
		m[i] = s
	}
	for i := 0; i < 100000000; i++ {
		delete(m, i)
	}
	for {
		debug.FreeOSMemory()
		time.Sleep(time.Second)
	}
}

// 按key排序
func TestSortMap(t *testing.T) {
	var m map[int]string
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Println("keys:", k, "values:", m[k])
	}
}

type data struct {
	name string
}

type printer interface {
	print()
}

func (p *data) print() {
	fmt.Println("name: ", p.name)
}

// 只要值是可寻址的，就可以在值上直接调用指针方法。即是对一个方法，它的 receiver 是指针就足矣。
// 但不是所有值都是可寻址的，比如 map 类型的元素、通过 interface 引用的变量：
func TestMapPointerReceiver(t *testing.T) {
	d1 := data{"one"}
	d1.print() // d1 变量可寻址，可直接调用指针 receiver 的方法

	//var in printer = data{"two"}
	//in.print() // 类型不匹配

	var in printer = &data{"two"}
	in.print() // 类型不匹配

	m := map[string]data{
		"x": data{"three"},
	}
	//m["x"].print() // m["x"] 是不可寻址的	// 变动频繁

	d := m["x"] // 编译通过
	d.print()

}
