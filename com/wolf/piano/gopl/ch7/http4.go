package main

import (
	"log"
	"net/http"
)

func main() {
	db := database3{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", http.HandlerFunc(db.list))
	http.HandleFunc("/price", http.HandlerFunc(db.price))
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
