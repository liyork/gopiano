package concurrent

import (
	"fmt"
	"sync"
	"testing"
)

type maps struct {
	mapinmap sync.Map
}

func TestSyncMapInMap(t *testing.T) {

	m := maps{}

	// 错误
	actualMap, ok := m.mapinmap.LoadOrStore("a", sync.Map{})
	// 正确1
	//actualMap, ok := m.mapinmap.LoadOrStore("a", &sync.Map{})
	fmt.Println("actualMap, ok:", actualMap, ok)

	// 错误
	// ./syncmap_test.go:20:35: cannot call pointer method on actualMap.(sync.Map)
	// 方法签名func (m *Map) LoadOrStore
	inmap, ok := actualMap.(sync.Map).LoadOrStore("b", "xxx1")
	// 正确2，单独声明一个变量，不报错。原因: 指针不能作为接收者, 需要需要定义变量来接送地址
	//i := actualMap.(sync.Map)
	//inmap, ok := i.LoadOrStore("b", "xxx1")
	// 正确1
	//inmap, ok := actualMap.(*sync.Map).LoadOrStore("b", "xxx1")
	fmt.Println("inmap, ok:", inmap, ok)

}
