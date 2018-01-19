package test

import (
	"bytes"
	"io"
	"os"
)

func CaptureStdout(f func()) string {
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	defer r.Close()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = stdOut

	var b bytes.Buffer
	io.Copy(&b, r)

	return b.String()
}
