package base

import "testing"

func sub1(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

func sub2(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

func sub3(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

func Add(i int, i2 int) int {
	return i + i2
}

// go test -v subunit_test.go
// go test subunit_test.go -v -run Sub/name1
func TestSub(t *testing.T) {

	t.Run("name1", sub1)
	t.Run("name2", sub2)
	t.Run("name3", sub3)
}
