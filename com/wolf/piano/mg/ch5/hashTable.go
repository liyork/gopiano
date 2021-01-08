package main

import "fmt"

const SIZE = 15

type Node struct {
	Value int
	Next  *Node
}

type HashTable struct {
	Table map[int]*Node
	Size  int
}

func hashFunction(i, size int) int {
	return i % size
}

func insert1(hash *HashTable, value int) int {
	index := hashFunction(value, hash.Size)
	// 新值位于槽的首部
	element := Node{Value: value, Next: hash.Table[index]}
	hash.Table[index] = &element
	return index
}

func traverse1(hash *HashTable) {
	for k := range hash.Table {
		if hash.Table[k] != nil {
			t := hash.Table[k]
			for t != nil {
				fmt.Printf("%d -> ", t.Value)
				t = t.Next
			}
		}
		fmt.Println()
	}
}

// go run hashTable.go
func main() {
	//baseTest()
	lookupTest()
}

func lookupTest() {
	table := make(map[int]*Node, SIZE)
	hash := &HashTable{Table: table, Size: SIZE}
	for i := 0; i < 120; i++ {
		insert1(hash, i)
	}
	fmt.Println(lookup(hash, 50))
	fmt.Println(lookup(hash, 121))
	fmt.Println(lookup(hash, 122))
	fmt.Println(lookup(hash, 123))
}

func baseTest() {
	table := make(map[int]*Node, SIZE)
	hash := &HashTable{Table: table, Size: SIZE}
	fmt.Println("Numbder of spaces:", hash.Size)
	for i := 0; i < 120; i++ {
		insert1(hash, i)
	}
	traverse1(hash)
}

func lookup(hash *HashTable, value int) bool {
	index := hashFunction(value, hash.Size)
	if hash.Table[index] != nil {
		t := hash.Table[index]
		for t != nil {
			if t.Value == value {
				return true
			}
			t = t.Next
		}
	}
	return false
}
