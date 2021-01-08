package main

import "fmt"

// go run method.go common.go
func main() {
	p := Point{1, 2}
	q := Point{4, 6}

	// method value，本身伴随方法接收者p
	distanceFromP := p.Distance
	fmt.Println(distanceFromP(q))
	var origin Point
	fmt.Println(distanceFromP(origin))

	//scaleP := p.ScaleBy
	//scaleP(2)
	//scaleP(3)
	//scaleP(10)

	// method expression,需要在调用时，第一个参数指定接受者
	var distance func(p, q Point) float64 = Point.Distance
	// 值类型
	fmt.Println(distance(p, q))
	fmt.Printf("%T\n", distance)
	fmt.Println(distance(q, p))

	var scale func(*Point, float64) = (*Point).ScaleBy
	// 指针类型
	scale(&p, 2)
	fmt.Println(p)
	fmt.Printf("%T\n", scale)
}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

type Path []Point

func (path Path) TranslateBy(offset Point, add bool) {
	// method expression
	var op func(p, q Point) Point
	if add {
		op = Point.Add
	} else {
		op = Point.Sub
	}
	for i := range path {
		path[i] = op(path[i], offset)
	}
}
