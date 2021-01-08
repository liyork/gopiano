package buffer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

// 多重读取器的接口
type MultipleReader interface {
	// 获取一个可关闭读取器的实例，持有该多重读取器中的数据
	Reader() io.ReadCloser
}

// 多重读取器的实现类型
type myMultipleReader struct {
	data []byte
}

// 新建并返回一个多重读取器的实例
func NewMultipleReader(reader io.Reader) (MultipleReader, error) {
	var data []byte
	var err error
	if reader != nil {
		data, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("multiple reader: couldn't create a new one: %s", err)
		}
	} else {
		data = []byte{}
	}
	return &myMultipleReader{
		data: data,
	}, nil
}

func (rr *myMultipleReader) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(rr.data))
}
