package main

import (
	"fmt"
	"net"
	"os"
)

// go run NSrecords.go mtsoukalos.eu
// go run NSrecords.go www.mtsoukalos.eu
// host -t ns www.mtsoukalos.eu
// host -t ns mtsoukalos.eu
func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Need a domain name!")
		return
	}

	domain := arguments[1]
	NSs, err := net.LookupNS(domain)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, NS := range NSs {
		fmt.Println(NS.Host)
	}
}
