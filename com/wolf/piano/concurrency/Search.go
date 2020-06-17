package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web1   = fakeSearch("web")
	Web2   = fakeSearch("web")
	Image1 = fakeSearch("image")
	Image2 = fakeSearch("image")
	Video1 = fakeSearch("video")
	Video2 = fakeSearch("video")
)

type Result2 string

type Search func(query string) Result2

func fakeSearch(kind string) Search {
	return func(query string) Result2 {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result2(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

func First2(query string, replicas ...Search) Result2 {
	c := make(chan Result2)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}

// 多并发请求
// Don't wait for slow servers. No locks. No condition variables. No callbacks.
func Google(query string) (results []Result2) {
	c := make(chan Result2)
	go func() { c <- First2(query, Web1, Web2) }()
	go func() { c <- First2(query, Image1, Image2) }()
	go func() { c <- First2(query, Video1, Video2) }()
	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
