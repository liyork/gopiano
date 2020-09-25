// +build hash !display_alternatives

package main

import "fmt"

type DisplayName string

func Print(name DisplayName) {
	fmt.Printf("%s\n", name)
}

func MakeDisplayName(name string) DisplayName {
	return DisplayName(name)
}
