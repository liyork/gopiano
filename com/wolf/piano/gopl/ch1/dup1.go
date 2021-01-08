package main

import (
	"bufio"
	"fmt"
	"os"
)

// go run dup1.go
// a
// b
// a
// mac ctrl+d
func main() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		//if input.Text() == "exit" {
		//	break
		//}
		counts[input.Text()]++
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
