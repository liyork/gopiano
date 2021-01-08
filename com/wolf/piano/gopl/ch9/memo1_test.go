package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

type Memo struct {
	f     Func
	cache map[string]result
}

type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

func (memo *Memo) Get(key string) (interface{}, error) {
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}

// go test -v memo1_test.go
func TestMemo1Base(t *testing.T) {
	m := New(httpGetBody)
	for _, url := range incomingURLs() {
		start := time.Now()
		value, err := m.Get(url)
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("%s, %s, %d bytes\n",
			url, time.Since(start), len(value.([]byte)))
	}
}

func incomingURLs() []string {
	return []string{"https://golang.org", "https://play.golang.org", "http://gopl.io"}
}

func httpGetBody(key string) (i interface{}, e error) {
	return nil, nil
}

// go test -run=TestMemo1Concurrent -race -v memo1_test.go
func TestMemo1Concurrent(t *testing.T) {
	m := New(httpGetBody)
	var n sync.WaitGroup
	for _, url := range incomingURLs() {
		n.Add(1)
		go func() {
			defer n.Done()
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s, %s, %d bytes\n",
				url, time.Since(start), len(value.([]byte)))
		}()
	}
	n.Wait()
}

type Memo2 struct {
	f     Func
	mu    sync.Mutex
	cache map[string]result
}

func (memo *Memo2) Get2(key string) (interface{}, error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	memo.mu.Unlock()
	return res.value, res.err
}

// 4 ,用锁+chan进行等待
type entry struct {
	res   result
	ready chan struct{}
}

type Memo4 struct {
	f     Func
	mu    sync.Mutex
	cache map[string]*entry
}

func (memo *Memo4) Get4(key string) (interface{}, error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)

		close(e.ready)
	} else {
		memo.mu.Unlock()

		<-e.ready
	}

	return e.res.value, e.res.err
}

// 5 单独服务进程处理所有请求
type request struct {
	key      string
	response chan<- result
}

type Memo5 struct {
	requests chan request
}

func New5(f Func) *Memo5 {
	memo := &Memo5{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo5) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo5) Close() {
	close(memo.requests)
}

// 进程中对每个请求，用两个协程交互，最后返回给用户结果
func (memo *Memo5) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	e.res.value, e.res.err = f(key)
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	<-e.ready
	response <- e.res
}
