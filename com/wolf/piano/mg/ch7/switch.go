package main

import "fmt"

type square1 struct {
	X float64
}

type circle1 struct {
	R float64
}

type rectangle struct {
	X float64
	Y float64
}

func tellInterface(x interface{}) {
	switch v := x.(type) {
	case square1:
		fmt.Println("This is a square!")
	case circle1:
		fmt.Printf("%v is a circle!\n", v)
	case rectangle:
		fmt.Println("This is a rectangle!")
	default:
		fmt.Printf("Unknown type %T!\n", v)
	}
}

// go run switch.go
func main() {
	x := circle1{R: 10}
	tellInterface(x)
	y := rectangle{X: 4, Y: 1}
	tellInterface(y)
	z := square1{X: 4}
	tellInterface(z)

	tellInterface(10)
}
