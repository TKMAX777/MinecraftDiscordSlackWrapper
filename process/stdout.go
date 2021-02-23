package process

import (
	"fmt"
	"io"
)

// StdoutReadCloser handles process's Stdout
//   this reader also exports its outputs to stdout
type StdoutReadCloser struct {
	readCloser io.ReadCloser
}

// NewStdout make stdin read and closer
func NewStdout(r io.ReadCloser) *StdoutReadCloser {
	return &StdoutReadCloser{r}
}

func (s StdoutReadCloser) Read(p []byte) (n int, err error) {
	n, err = s.readCloser.Read(p)
	fmt.Printf("%s", p[:n])

	if err != nil {
		return
	}
	return
}

// Close ...
func (s StdoutReadCloser) Close() error {
	return s.readCloser.Close()
}
