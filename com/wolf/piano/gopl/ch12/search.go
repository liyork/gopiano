package main

import (
	"fmt"
	"net/http"
)

// go run search.go
// ./fetch 'http://localhost:12345/search'
// ./fetch 'http://localhost:12345/search?l=golang&l=programming'
// ./fetch 'http://localhost:12345/search?l=golang&l=programming&max=100'
// ./fetch 'http://localhost:12345/search?x=true&l=golang&l=programming'
// ./fetch 'http://localhost:12345/search?q=hello&x=123'

func search(resp http.ResponseWriter, req *http.Request) {
	var data struct {
		Labels     []string `http:"1"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}
	data.MaxResults = 10
	if err := Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	// handler
	fmt.Fprintf(resp, "Search: %+v\n", data)
}
