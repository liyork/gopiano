package base

import (
	"fmt"
	"testing"
)

func TestFmtBase(t *testing.T) {

	sprintf := fmt.Sprintf("[%d, %d, %d]", 1, 2, 3)
	fmt.Println(sprintf)
}
