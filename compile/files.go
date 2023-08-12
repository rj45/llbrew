package compile

import (
	"io"
	"log"
	"os"
)

func createFile(filename string) io.WriteCloser {
	if filename == "-" {
		return nopWriteCloser{os.Stdout}
	}
	if filename == "" {
		return nopWriteCloser{}
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

type nopWriteCloser struct{ w io.Writer }

func (nopWriteCloser) Close() error {
	return nil
}

func (n nopWriteCloser) Write(p []byte) (int, error) {
	if n.w == nil {
		return len(p), nil
	}
	return n.w.Write(p)
}
