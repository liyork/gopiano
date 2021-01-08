package main

import (
	"fmt"
	"os"
	"strconv"
)

// go build cf.go tempconv.go
// ./cf 32
// ./cf 212
// ./cf -40

func main() {
	for _, arg := range os.Args[1:] {
		t, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cf: %v\n", err)
			os.Exit(1)
		}
		// 强转
		f := Fahrenheit(t)
		c := Celsius(t)
		fmt.Printf("%s = %s, %s = %s\n",
			f, FToC(f), c, CToF(c))
	}
}
