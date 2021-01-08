package documentMe

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/mg/ch11"
)

func ExampleS1() {
	fmt.Println(main.S1("123456789"))
	fmt.Println(main.S1(""))
	// Output:
	// 9
	// 0
}

// godoc -http=":8080"
// http://localhost:8080/pkg/
func ExampleF1() {
	fmt.Println(main.F1(10))
	fmt.Println(main.F1(2))
	// Output:
	// 1
	// 55
}
