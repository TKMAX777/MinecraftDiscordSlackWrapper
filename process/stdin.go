package process

import (
	"fmt"
	"io"
	"os"
)

// StdinWriteCloser wraps io.WriteCloser
//
//	start method starts synchronizing process's stdin with stdin
type StdinWriteCloser struct {
	writeCloser io.WriteCloser
}

// NewStdin make stdin write and closer
func NewStdin(w io.WriteCloser) *StdinWriteCloser {
	var stdin = &StdinWriteCloser{w}
	go stdin.start()
	return stdin
}

func (w StdinWriteCloser) start() {
	var err error
	var buf = make([]byte, 100)
	var n, m, i int

	for {
		n, err = os.Stdin.Read(buf)
		switch err {
		case nil:
		case io.EOF:
			return
		default:
			fmt.Printf("(Input)Error:%s\n", err.Error())
			continue
		}

		m = 0
		for m < n {
			i, err = w.writeCloser.Write(buf[m:n])
			if err != nil {
				fmt.Printf("(Input)Error:%s\n", err.Error())
				return
			}
			m += i
		}
	}
}

func (w StdinWriteCloser) Write(p []byte) (n int, err error) {
	return w.writeCloser.Write(p)
}

// Close ...
func (w StdinWriteCloser) Close() error {
	return w.writeCloser.Close()
}
