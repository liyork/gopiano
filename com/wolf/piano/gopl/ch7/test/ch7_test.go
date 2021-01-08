package test

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestBase(t *testing.T) {
	var w io.Writer
	w = os.Stdout
	w = new(bytes.Buffer)
	//w = time.Second

	var rwc io.ReadWriteCloser
	rwc = os.Stdout
	//rwc = new(bytes.Buffer)

	w = rwc
	//rwc = w

	// interface value
	w = os.Stdout
	w = new(bytes.Buffer)
	w = nil
}
