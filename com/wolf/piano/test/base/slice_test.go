package base

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

//切片的本质,是一种新定义的一种数据结构
//type slice struct {
//            *Pointer
//            len
//            cap
//        }

//数组类型是很有用的，但是不太灵活，所以Go代码中很少看到它们。但是切片类型却是很常见的，因为它基于数组类型提供了强大的功能和开发便利。
//切片类型的定义如[]T，其中T是切片中元素的类型。与数组类型不同，切片类型没有固定的长度。
//定义一个切片和定义一个数组的语法相似，唯一的不同是不需要定义切片长度。
func TestSliceBase(t *testing.T) {
	letters := []string{"a", "b", "c", "d"}
	fmt.Println(letters)

	// make分配一个数组，并且返回一个指向该数组的切片。
	var s = make([]byte, 5, 5)
	fmt.Println(s)
	//如果没有传入cap参数，它的默认值是传入的长度。这是上面代码的一个简洁版本。
	s = make([]byte, 5)
	fmt.Println(s)
	fmt.Println(len(s), cap(s))

	//切片的零值为nil
	var snil []int
	fmt.Println(len(snil), cap(snil)) // 0, 0

	// 可以通过“切”一个数组或者是切片，来生成新的切片。
	//b[1:4]会返回一个新的切片，包含的元素为b中的第1到第4-1的元素(最后不包含)
	b := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	// 从0开始，不包括endIndex
	bs := b[1:4] // bs和b中的元素占用同一块内存
	fmt.Printf("bs:%s\n", bs)

	//基于数组创建切片
	x := [3]string{"Лайка", "Белка", "Стрелка"}
	s1 := x[:] // s1为指向x切片的引用
	fmt.Println(s1)
}

// 引用相同底层数组
func TestSliceUpdate(t *testing.T) {
	a := []int{1, 2, 3}
	b := a[0:1]
	b[0] = 44
	fmt.Println(a, " xxxx ", b)
}

// slice没有删除元素的函数。可以使用数组切片重新组合的方式来删除一个或多个项。不过从数组切片这种数据结构来看，本身并不适合做删除操作，所以尽量减少使用
func TestSliceDelete(t *testing.T) {
	a := []int{1, 2, 3}
	index := 1
	// 删除下标index的元素
	a = append(a[:index], a[index+1:]...)
	fmt.Println(a)
}

// slice和数组在声明时的区别：
// 声明数组时，方括号内写明了数组的长度或使用...自动计算长度，
// 声明slice时，方括号内没有任何字符。
func TestSliceArrayDiff(t *testing.T) {
	vs := []interface{}{
		[]int{},            // slice 切片
		[]int{1, 2, 3},     // slice 切片
		[]int{1, 2, 3}[:],  //切片再切还是切片
		make([]int, 3, 10), //标准的slice 定义方式
		[3]int{1, 2, 3},    //array 数组，确定数组长度
		[...]int{1, 2, 3},  //array 数组，由编译器自动计算数组长度。
	}
	for i, v := range vs {
		rv := reflect.ValueOf(v)
		fmt.Println(i, rv.Kind())
	}
}

// 按照目标的len进行拷贝
func TestSliceCopy(t *testing.T) {
	var source = make([]string, 0)
	for i := 0; i < 10; i++ {
		source = append(source, fmt.Sprintf("%v", i))

	}

	// 按照len最小的拷贝
	var destination = make([]string, 0, 10)
	copyLen := copy(destination, source)
	fmt.Printf("copy to destination(len=%d)\t%v\n", len(destination), destination)

	destination = make([]string, 5)
	copyLen = copy(destination, source)
	fmt.Printf("copy to destination(len=%d)\tcopylen=%d\t%v\n", len(destination), copyLen, destination)

	destination = make([]string, 10)
	copyLen = copy(destination, source)
	fmt.Printf("copy to destination(len=%d)\tcopylen=%d\t%v\n", len(destination), copyLen, destination)
}

//切片是数组段的描述符。它包含了一个指向数组的指针ptr，数据段的长度len和容量cap
//长度是切片指向内容中元素的个数。容量是底层数组中的元素个数（从切片指向的元素开始计数）
//切片操作并不会拷贝s中的数据，而是创建一个新的切片指向原来的数组，这让切片操作就像操作数组索引一样高效。
// 因此，对切片的元素进行修改，会修改原始切片的元素。
// 要增加切片的容量，必须新建一个容量更大的切片，然后将之前的切片的数据拷贝到新的切片中
//对一个切片进行切片不会拷贝切片指向的数组。这个数组会一致保存在内存中，直到不再被引用。
// 有时这样会导致程序会将所有的数据保存在内存中，即使只有一小部分数据是被需要的。
//只要这个切片一直保留着，垃圾回收将不能释放保存了所有数据的数组。文件一小部分有用的数据将会让所有的数据一直保存在内存中
//要解决这个问题，可以先将有用的数据先保存到一个新的切片，然后返回新的切片。看来回收是依据引用，那么没被引用的数组元素会被回收？
func TestSliceInter(t *testing.T) {
	d := []byte{'r', 'o', 'a', 'd'}
	e := d[2:]
	// e == []byte{'a', 'd'}
	e[1] = 'm'
	// e == []byte{'a', 'm'}
	// d == []byte{'r', 'o', 'a', 'm'}
	fmt.Printf("%s\n", e)

	//新建一个容量是s两倍的切片t，然后将s的数据拷贝到t中，最后将t赋值给s:
	q := make([]byte, len(d), (cap(d)+1)*2) // +1对应 cap(s) == 0的情况
	for i := range d {
		q[i] = d[i]
	}
	fmt.Println(len(d), cap(d))
	//d = q
	//fmt.Println(len(d),cap(d))

	//使用内置的copy函数可以简化上面的代码。顾名思义，copy将数据从一个切片拷贝到另一个切片，并返回拷贝元素的数量
	q1 := make([]byte, len(d), (cap(d)+1)*2)
	fmt.Println(len(d), cap(d))
	copy(q1, d)
	d = q1
	fmt.Println(len(d), cap(d))

	//一个常见的操作是在切片的末尾添加一个元素。下面的函数在一个切片的末尾增加一个元素，在容量不够的情况下增加切片的容量，并且返回更新后的切片
	//函数append 将x添加到s末尾，如果需要就扩展s的容量。
	a := make([]int, 1)
	// a == []int{0}
	a = append(a, 1, 2, 3)
	// a == []int{0, 1, 2, 3}

	//使用...将一个切片添加到另外一个切片末尾
	a1 := []string{"John", "Paul"}
	b1 := []string{"George", "Ringo", "Pete"}
	a1 = append(a1, b1...) // 等同于append(a, b[0], b[1], b[2])
	//  a == []string{"John", "Paul", "George", "Ringo", "Pete"}
}

//append()函数,在调用函数时，在栈区里面把1 2 3 添加到a里面然后重新分配了地址，而原来的s切片还是指向原来地址，所以需要重新指向然后返回
// 切片用append函数的时候一定要注意，因为如果容量不足的时候会自动扩充，
// 如果原来的地址后面没有足够的空间那么就会重新找一个足够大的空间来储存，所以切片利用append的时候地址是有可能变化的。
// 扩容方式与vector类似，1024字节以下是每次cap增加一倍，1024以上是每次cap增加1/4
func TestSliceAppend(t *testing.T) {
	var s []int = []int{89, 4, 5, 6}
	sp := &s
	fmt.Printf("sAddr:%p,sPointAddr:%p\n", s, sp)
	for i := 1; i < 100; i++ {
		*sp = append(*sp, i*200000)
	}
	// 变量s地址变化，但是sp指针未变化
	fmt.Printf("sAddr:%p,sPointAddr:%p\n", s, sp)
}

func TestSliceSort(t *testing.T) {
	s := []string{"b", "d", "c"}
	sort.Strings(s)

	fmt.Println(s)
}

func TestConvert(t *testing.T) {
	var buf [1024]byte
	fmt.Println("xx", reflect.TypeOf(buf).Kind())
	fmt.Println("xx", reflect.TypeOf(buf[:]).Kind())
	fmt.Println("xx", len(buf))
	fmt.Println("xx", len(buf[:]))

	var buf1 = make([]byte, 3333444433)
	fmt.Println("xx", len(buf1[:333333]))
	fmt.Println("xx", len(buf1))
}

// golang数组切片解决的问题：golang数组长度在定义之后无法再次修改，并且数组是值类型，每次传递都将产生一份副本
//golang数组切片拥有独立的数据结构，可抽象为3个变量：一个指向原数组的指针，数组切片中元素个数，数组切片分配的存储空间
func TestSliceFromSlice(t *testing.T) {
	oldSlice := make([]int, 5, 10)
	// 范围不能超过oldSlice的capacity
	newSlice := oldSlice[:4]

	fmt.Println("Length   of oldSlice: ", len(oldSlice)) // 5
	fmt.Println("Capacity of oldSlice: ", cap(oldSlice)) // 10
	fmt.Println("oldSlice: ", oldSlice)

	fmt.Println("Length   of newSlice: ", len(newSlice)) // 8
	// 相同capacity
	fmt.Println("Capacity of newSlice: ", cap(newSlice)) // 10
	fmt.Println("newSlice: ", newSlice)

}

// 与数组相比，数组切片多了一个capacity的概念，即当前容纳的元素个数len和分配的空间capacity可以是两个不同的值。
//capacity，可以理解为最大容纳元素个数，最大容纳元素个数减去当前容纳元素个数剩下的空间是隐藏的，不能直接使用。如果要往隐藏空间中新增元素，可以使用append()函数。
//取得当前容纳元素个数可以使用len()函数，取得最大容纳元素个数可以使用cap()函数。

func TestEstimate(t *testing.T) {
	b := make([]int, 1024)
	b[0] = 1
	// 导致扩容
	//b = append(b, 99)
	fmt.Println("len:", len(b), "cap:", cap(b))
}

// 扩容
//两条规则：
//如果切片的容量小于1024个元素，那么扩容的时候slice的cap就翻番，乘以2；一旦元素个数超过1024个元素，增长因子就变成1.25，即每次增加原来容量的四分之一。
//如果扩容之后，还没有触及原数组的容量，那么，切片中的指针指向的位置，就还是原数组，如果扩容之后，超过了原数组的容量，那么，Go就会开辟一块新的内存，把原来的值拷贝过来，这种情况丝毫不会影响到原数组。
func TestScale(t *testing.T) {
	arr := [4]int{10, 20, 30, 40}
	slice := arr[0:2]
	fmt.Println("slice==>", slice, len(slice), cap(slice))
	testSlice1 := slice // [10,20]
	fmt.Println("testSlice1==>", testSlice1, len(testSlice1), cap(testSlice1))
	testSlice2 := append(append(append(slice, 1), 2), 3) // [10,20,1,2,扩容*2-->,3,0,0,0]
	fmt.Println("testSlice2==>", testSlice2, len(testSlice2), cap(testSlice2))
	slice[0] = 11

	fmt.Println(testSlice1[0])
	fmt.Println(testSlice2[0])
}
