package process

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Process handles process controls
type Process struct {
	cmd  *exec.Cmd
	pipe pipe
}

// Pipes returns stdin / stdout / stderr pipes
func (p *Process) Pipes() (io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	return p.pipe.Pipes()
}

// Interrupt sends keyboard interrupt to the process and wait until it stops
func (p *Process) Interrupt() {
	p.pipe.Close()

	p.cmd.Process.Signal(os.Interrupt)
	p.cmd.Wait()

	return
}

// Kill kills the process immediately
func (p *Process) Kill() {
	p.pipe.Close()
	p.cmd.Process.Kill()
	return
}

// Start excutes the commmand
func (p *Process) Start(name string, arg ...string) error {
	var err error

	if p.cmd != nil {
		return fmt.Errorf("Already cmd started")
	}

	p.cmd = exec.Command(name, arg...)

	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Failed to build pipe\n%s", err.Error())
	}
	p.pipe.stdOut = NewStdout(stdout)

	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Failed to build pipe\n%s", err.Error())
	}
	p.pipe.stdErr = NewStderr(stderr)

	stdin, err := p.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("Failed to build pipe\n%s", err.Error())
	}
	p.pipe.stdIn = NewStdin(stdin)

	return p.cmd.Start()
}
