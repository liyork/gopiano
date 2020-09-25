package main

import (
	"bytes"
	"log"
	"math/rand"
	"time"
)

func test() {

	log.Println(" ===> loop begin.")
	for i := 0; i < 10; i++ {
		log.Println(genSomeBytes())
	}

	log.Println(" ===> loop end.")
}

func genSomeBytes() *bytes.Buffer {

	var buff bytes.Buffer

	for i := 1; i < 20; i++ {
		buff.Write([]byte{'0' + byte(rand.Intn(10))})
	}

	return &buff
}

// time go run timetest.go
// /usr/bin/time -v go run timetest.go(mac 不用加-v)
func main() {

	for i := 0; i < 2; i++ {
		test()
		time.Sleep(time.Second * 1)
	}

}
