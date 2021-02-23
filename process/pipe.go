package process

import (
	"fmt"
	"io"
)

type pipe struct {
	isPipeExported bool

	stdIn  *StdinWriteCloser
	stdOut *StdoutReadCloser
	stdErr *StderrReadCloser
}

// Pipes returns stdin / stdout / stderr pipes
func (p *pipe) Pipes() (io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	if p.isPipeExported {
		return nil, nil, nil, fmt.Errorf("The pipes already exported")
	}
	return p.stdIn, p.stdOut, p.stdErr, nil
}

// Close closes all pipes
func (p *pipe) Close() {
	p.stdIn.Close()
	p.stdOut.Close()
	p.stdErr.Close()

	return
}
