package base

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/test/protobuf/other"
	"testing"
)

func TestImport(t *testing.T) {
	response := SearchResponse{Results: []*other.Result3{}}
	fmt.Println("response：", response)
}
