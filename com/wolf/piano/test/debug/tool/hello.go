package tool

import "fmt"

func main() {
	s := test("qqq,")
	fmt.Println("s:", s)
}

func test(a string) string {
	return a + "xx1"
}
