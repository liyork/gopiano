package main

// 两个文件中一个是hash !display_alternatives tag,或者关系。一个是int tag，
// 需要引入第三个tag(display_alternatives)来区分编译的文件。否则，只要不带！的tag都会被编译进包。

// go build -tags int  --报错
// go build -tags "display_alternatives int"  --编译int
// go build -tags hash  --编译hash
// ./tags
func main() {
	var name DisplayName
	name = MakeDisplayName("FAD9C812")
	println(name)
}
