package main

import (
	"fmt"
	"testing"
)

// go test
// 测试Palindrome，若返回false则错误
func TestPalindrome(t *testing.T) {
	if !IsPalindrome("detartrated") {
		t.Errorf(`IsPalindrome("detartrated") = false`)
	}
	if !IsPalindrome("kayak") {
		t.Errorf(`IsPalindrome("kayak") = false`)
	}
}

// 测试非Palindrome，若true则错误
func TestNonPalindrome(t *testing.T) {
	if IsPalindrome("palindrome") {
		t.Errorf(`IsPalindrome("palindrome") = true`)
	}
}

// go test -v
// go test -v -run="TestCanalPalindrome"
func TestCanalPalindrome(t *testing.T) {
	input := "A man, a plan, canal: Panama"
	if !IsPalindrome(input) {
		t.Errorf(`IsPalindrome("%q") = false`, input)
	}
}

// go test -v -run="TestIsPalindrome"
func TestIsPalindrome(t *testing.T) {
	var tests = []struct {
		intput string
		want   bool
	}{
		{"", true},
		{"a", true},
		{"aa", true},
		{"ab", false},
		{"kayak", true},
		{"detartrated", true},
	}
	for _, test := range tests {
		if got := IsPalindrome2(test.intput); got != test.want {
			t.Errorf("IsPalindrome2(%q) = %v", test.intput, got)
		}
	}
}

// go test -v -run=Coveragge eval
// go tool cover
// go test -run=Coverage -coverprofile=c.out eval
// go tool cover -html=c.out

// go test -bench=.
// go test -bench=. -benchmem
func BenchmarkIsPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPalindrome2("A man, a plan, canal: Panama")
	}
}

// go test -cpuprofile=cpu.out
// go test -blockprofile=block.out
// go test -memprofile=mem.out
// go test -run=NONE -bench=ClientServerParallelTLS64 -cpuprofile=cpu.log net/http
// go tool pprof -test -nodecount=10 ./http.test cpu.log

func ExampleIsPalindrome() {
	fmt.Println(IsPalindrome("A man, a plan, acanal: Panama"))
	fmt.Println(IsPalindrome("palindrome"))
	// Output:
	// true
	// false
}
