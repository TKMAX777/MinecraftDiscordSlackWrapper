package process

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// StdinWriteCloser wraps io.WriteCloser
//    start method starts synchronizing process's stdin with stdin
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
	var input []byte
	var err error

	for {
		input, err = w.scan()
		if err != nil {
			fmt.Printf("(Input)Error:%s\n", err.Error())
			continue
		}

		_, err = w.Write(append(input, []byte("\n")...))
		if err != nil {
			fmt.Printf("(Input)Error:%s\n", err.Error())
			continue
		}
	}
}

func (w StdinWriteCloser) scan() ([]byte, error) {
	var rdr = bufio.NewReaderSize(os.Stdin, bufio.MaxScanTokenSize)
	var buf []byte = []byte{}

	for {
		l, p, e := rdr.ReadLine()
		if e != nil {
			return nil, e
		}
		buf = append(buf, l...)

		if !p {
			break
		}
	}
	return buf, nil
}

func (w StdinWriteCloser) Write(p []byte) (n int, err error) {
	return w.writeCloser.Write(p)
}

// Close ...
func (w StdinWriteCloser) Close() error {
	return w.writeCloser.Close()
}
