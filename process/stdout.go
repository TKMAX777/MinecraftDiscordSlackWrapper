package process

import (
	"io"
	"os"
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
	os.Stderr.Write(p[:n])

	return
}

// Close ...
func (s StdoutReadCloser) Close() error {
	return s.readCloser.Close()
}
