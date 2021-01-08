package main

import "fmt"

// 后进先出

type Node struct {
	Value int
	Next  *Node
}

var size = 0
var queue = new(Node)

// 链表首部添加
func Push(v int) bool {
	if queue == nil {
		queue = &Node{v, nil}
		size = 1
		return true
	}

	temp := &Node{v, nil}
	temp.Next = queue
	queue = temp
	size++
	return true
}

// 从头部出
func Pop(t *Node) (int, bool) {
	if size == 0 {
		return 0, false
	}

	if size == 1 {
		size = 0
		queue = nil
		return t.Value, true
	}

	queue = queue.Next
	size--
	return t.Value, true
}

func traverse(t *Node) {
	if size == 0 {
		fmt.Println("Empty Stack!")
		return
	}

	for t != nil {
		fmt.Printf("%d -> ", t.Value)
		t = t.Next
	}
	fmt.Println()
}

func main() {
	queue = nil
	v, b := Pop(queue)
	if b {
		fmt.Print(v, " ")
	} else {
		fmt.Println("Pop() failed!")
	}

	Push(100)
	traverse(queue)
	Push(200)
	traverse(queue)

	for i := 0; i < 10; i++ {
		Push(i)
	}

	for i := 0; i < 15; i++ {
		v, b := Pop(queue)
		if b {
			fmt.Print(v, " ")
		} else {
			break
		}
	}
	fmt.Println()
	traverse(queue)
}
