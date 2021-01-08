package main

import (
	"flag"
	"fmt"
)

// go run simpleFlag.go -O 100
// go run simpleFlag.go -O=100
// go run simpleFlag.go -O=100 -k
// go run simpleFlag.go -O=100 -k false
// go run simpleFlag.go -O=100 -k=false
// go run simpleFlag.go -O=notAnInterger
func main() {
	minusK := flag.Bool("k", true, "k")
	minusO := flag.Int("O", 1, "O")
	flag.Parse()

	valueK := *minusK
	valueO := *minusO
	valueO++

	fmt.Println("-k:", valueK)
	fmt.Println("-O:", valueO)
}
