package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"
)

func main() {
	//toStruct()
	//toInterface()
	toJsonString()
}

// json的key是Foo，匹配过程：
// 首先查找tag含有Foo的可导出的struct字段(首字母大写)
// 其次查找字段名是Foo的导出字段
// 最后查找类似FOO或者FoO这样的除了首字母之外其他大小写不敏感的导出字段
// 注意：能够被赋值的字段必须是可导出字段（**即首字母大写**）。同时JSON解析的时候只会解析能找到的字段，如果找不到的字段会被忽略
func toStruct() {

	type Server struct {
		ServerName string
		ServerIP   string
	}

	type Serverslice struct {
		Servers []Server
	}

	var s Serverslice
	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},
	            {"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(str), &s)
	fmt.Println(s)
	fmt.Println(s.Servers[0].ServerIP)
}

// Go类型和JSON类型的对应关系如下：
//  bool代表JSON booleans
//  float64代表JSON numbers
//  string代表JSON strings
//  nil 代表JSON null
func toInterface() {
	b := []byte(`{"Name":"Wednesday", "Age":6, "Parents": [ "Gomez", "Moticia" ]}`)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println(err)
	}

	m := f.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}

//  在定义struct tag的时候需要注意几点：
//  字段的tag是"-"，那么这个字段不会输出到JSON
//  tag中带有自定义名称，那么这个自定义名称会出现在JSON的字段名中，例如xxxx
//  tag中如果带有“omitempty”选项，那么如果该字段值为空，就不会输出到JSON串中，默认是输出默认值
//  如果字段类型是bool,string,int,int64等，而tag中带有“,string”选项，那么这个字段在输出到JSON的时候会把该字段对应的值转换成JSON字符串，即\"xxx\"这种
// JSON对象只支持string作为key,所以要编码一个map,那么必须是map[string]T这种类型
// Channel,complex和function是不能被编码成JSON的
// 嵌套的数据时不能编码的，不然会让JSON编码进入死循环
// 指针在编码的时候会输出指针指向的内容，而空指针会输出null
func toJsonString() {
	type Server struct {
		//ID不会导出到JSON
		ID         int    `json:"-"` // 要转换必须大写，类型也得匹配
		ServerName string `json:"serverName,string"`
		ServerIp   string `json:"xxxx,omitempty"`
	}

	type Serverslice struct {
		Servers []Server `json:"servers"`
	}

	var s Serverslice
	s.Servers = append(s.Servers, Server{ServerName: "servername1", ServerIp: "ip1"})
	s.Servers = append(s.Servers, Server{ServerName: "servername2"})

	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err: ", err)
	}

	fmt.Println(string(b))
}

type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color, omitempty"`
	Actors []string
}

func Test_convert(t *testing.T) {

	var movies = []Movie{
		{Title: "Casabanca", Year: 1942, Color: false,
			Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
		{Title: "Casabanca2", Year: 1962, Color: true,
			Actors: []string{"Humphrey Paul"}},
	}

	data, err := json.Marshal(movies)
	if err != nil {
		log.Fatalf("Json marshaling failed：%s", err)
	}
	fmt.Printf("%s\n", data)
	var movies2 []struct{ Title string }
	// movies2 := make([]Movie, 10)
	if err3 := json.Unmarshal(data, &movies2); err3 != nil {
		log.Fatalf("JSON unmarshling failed: %s", err)
	}
	fmt.Println("movies2*****************", movies2)
	data2, err2 := json.MarshalIndent(movies, "", " ")
	if err2 != nil {
		log.Fatalf("Json marshlindent failed：%s", err)
	}
	fmt.Printf("%s\n", data2)
}

type Stu struct {
	Name  string `json:"name"`
	Age   int
	HIgh  bool
	sex   string
	Class *Class `json:"class"`
}

type Class struct {
	Name  string
	Grade int
	//C     *Class//循环的数据结构也不行，它会导致marshal陷入死循环。
}

// 只要是可导出成员（变量首字母大写），都可以转成json。因成员变量sex是不可导出的，故无法转成json。
// 如果变量打上了json标签，如Name旁边的 `json:"name"` ，那么转化成的json key就用该标签“name”，否则取变量名作为key，如“Age”，“HIgh”。
// Channel， complex 以及函数不能被编码json字符串,
// 指针变量，编码时自动转换为它所指向的值，如cla变量。
//（当然，不传指针，Stu struct的成员Class如果换成Class struct类型，效果也是一模一样的。只不过指针更快，且能节省内存空间。）
// json编码成字符串后就是纯粹的字符串了
// 无论是string，int，bool，还是指针类型等，都可赋值给interface{}类型，且正常编码
func TestMarshal(t *testing.T) {
	stu := Stu{
		Name: "张三",
		Age:  18,
		HIgh: true,
		sex:  "男",
	}

	//指针变量
	cla := new(Class)
	cla.Name = "1班"
	cla.Grade = 3
	//cla.C = cla
	stu.Class = cla

	jsonStu, err := json.Marshal(stu)
	if err != nil {
		fmt.Println("生成json字符串错误")
	}

	//jsonStu是[]byte类型，转化成string类型便于查看
	fmt.Println(string(jsonStu))
}

type StuRead struct {
	Name  interface{} `json:"name"`
	Age   interface{}
	HIgh  interface{}
	sex   interface{}
	Class interface{} `json:"class"`
	//Class Class `json:"class"`
	//Class *Class `json:"class"`
	Test interface{}
}

type StuReadCommon struct {
	Name  interface{} `json:"name"`
	Age   interface{}
	HIgh  interface{}
	sex   interface{}
	Class interface{} `json:"class"`
	//Class Class `json:"class"`
	//Class *Class `json:"class"`
	Test interface{}
}

type StuRead3 struct {
	*StuReadCommon
}

// 解析时，接收体可自行定义。json串中的key自动在接收体中寻找匹配的项进行赋值。匹配规则是：
//
//先查找与key一样的json标签，找到则赋值给该标签对应的变量(如Name)。
//没有json标签的，就从上往下依次查找变量名与key一样的变量，如Age。或者变量名忽略大小写后与key一样的变量。如HIgh，Class。第一个匹配的就赋值，后面就算有匹配的也忽略。
//(前提是该变量必需是可导出的，即首字母大写)。
// 不可导出的变量无法被解析（如sex变量，虽然json串中有key为sex的k-v，解析后其值仍为nil,即空值）
// 当接收体中存在json串中匹配不了的项时，解析会自动忽略该项，该项仍保留原值。如变量Test，保留空值nil。
// 变量Class貌似没有解析为我们期待样子。因为此时的Class是个interface{}类型的变量，而json串中key为CLASS的value是个复合结构，不是可以直接解析的简单类型数据（如“张三”，18，true等）。所以解析时，由于没有指定变量Class的具体类型，json自动将value为复合结构的数据解析为map[string]interface{}类型的项。也就是说，此时的struct Class对象与StuRead中的Class变量没有半毛钱关系，故与这次的json解析没有半毛钱关系
func TestUnMarshal(t *testing.T) {
	//json字符中的"引号，需用\进行转义，否则编译出错
	//json字符串沿用上面的结果，但对key进行了大小的修改，并添加了sex数据
	data := "{\"name\":\"张三\",\"Age\":18,\"high\":true,\"sex\":\"男\",\"CLASS\":{\"naME\":\"1班\",\"GradE\":3}}"
	str := []byte(data)

	//第二个参数必须是指针，否则无法接收解析的数据，如stu仍为空对象StuRead{}
	stu := &StuRead{}
	printType(stu)
	err := json.Unmarshal(str, stu)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(stu)

	fmt.Println("-----------")
	printType(stu)
}

//如果不想指定Class变量为具体的类型，仍想保留interface{}类型，但又希望该变量可以解析到struct Class对象中，这时候该怎么办呢？
type StuRead2 struct {
	Name  interface{}
	Age   interface{}
	HIgh  interface{}
	Class json.RawMessage `json:"class"` //注意这里
}

// 在第一次json解析时，变量Class的类型是json.RawMessage。对该变量进行二次json解析，我们只需再定义一个新的接受体即可，如json.Unmarshal(stu.Class,cla)
func TestUnMarshal2(t *testing.T) {
	data := "{\"name\":\"张三\",\"Age\":18,\"high\":true,\"sex\":\"男\",\"CLASS\":{\"naME\":\"1班\",\"GradE\":3}}"
	str := []byte(data)
	stu := StuRead2{}
	err := json.Unmarshal(str, &stu)
	if err != nil {
		fmt.Println(err)
	}

	//注意这里：二次解析！
	cla := new(Class)
	json.Unmarshal(stu.Class, cla)

	fmt.Println("stu:", stu)
	fmt.Println("string(stu.Class):", string(stu.Class))
	fmt.Println("class:", cla)
}

// 通过继承属性方式，不要再指定StuReadCommon的`json
func TestUnMarshal3(t *testing.T) {
	data := "{\"name\":\"张三\",\"Age\":18,\"high\":true,\"sex\":\"男\",\"CLASS\":{\"naME\":\"1班\",\"GradE\":3}}"
	str := []byte(data)

	//第二个参数必须是指针，否则无法接收解析的数据，如stu仍为空对象StuRead{}
	stu := &StuRead3{}
	json.Unmarshal(str, stu)

	fmt.Printf("stu:,%+v", stu)
}

// interface{}类型变量在json解析前，打印出的类型都为nil，就是没有具体类型，这是空接口（interface{}类型）的特点。
//
//json解析后，json串中value，只要是”简单数据”，都会按照默认的类型赋值，如”张三”被赋值成string类型到Name变量中，数字18对应float64，true对应bool类型。
// “简单数据”：是指不能再进行二次json解析的数据，如”name”:”张三”只能进行一次json解析。
//“复合数据”：类似”CLASS\”:{\”naME\”:\”1班\”,\”GradE\”:3}这样的数据，是可进行二次甚至多次json解析的，因为它的value也是个可被解析的独立json。即第一次解析key为CLASS的value，第二次解析value中的key为naME和GradE的value
// 对于”复合数据”，如果接收体中配的项被声明为interface{}类型，go都会默认解析成map[string]interface{}类型。如果我们想直接解析到struct Class对象中，可以将接受体对应的项定义为该struct类型 Class *Class `json:"class"`
func printType(stu *StuRead) {
	nameType := reflect.TypeOf(stu.Name)
	ageType := reflect.TypeOf(stu.Age)
	highType := reflect.TypeOf(stu.HIgh)
	sexType := reflect.TypeOf(stu.sex)
	classType := reflect.TypeOf(stu.Class)
	testType := reflect.TypeOf(stu.Test)

	fmt.Println("nameType:", nameType)
	fmt.Println("ageType:", ageType)
	fmt.Println("highType:", highType)
	fmt.Println("sexType:", sexType)
	fmt.Println("classType:", classType)
	fmt.Println("testType:", testType)
}

type Movie2 struct {
	Title  string
	Actors []string
}

type Movie3 struct {
	Title1 string
	m2     Movie2
}

// 对于m2结构体，反序列化后是有默认结构体，里面属性是默认值
func TestUnMarshalNilStruct(t *testing.T) {
	data := "{\"title1\":\"张三\"}"
	str := []byte(data)

	stu := &Movie3{}
	err := json.Unmarshal(str, stu)

	fmt.Println(stu, err)
}
