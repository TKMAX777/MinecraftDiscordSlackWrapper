package process

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessWorkingDirectory(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	p := Process{WorkingDirectory: directory}
	if err := p.Start(executable, "-test.run=^TestProcessDirectoryHelper$", "--", directory); err != nil {
		t.Fatal(err)
	}
	stdin, stdout, stderr, err := p.Pipes()
	if err != nil {
		t.Fatal(err)
	}
	defer stdout.Close()
	defer stderr.Close()
	stdin.Close()
	output, err := io.ReadAll(stdout)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	if string(output) != "directory-ok" {
		t.Fatalf("subprocess did not use jar directory: %s", output)
	}
}

func TestProcessDirectoryHelper(t *testing.T) {
	if len(os.Args) < 3 || os.Args[len(os.Args)-2] != "--" {
		return
	}
	actual, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}
	actualInfo, err := os.Stat(filepath.Clean(actual))
	if err != nil {
		os.Exit(1)
	}
	wantedInfo, err := os.Stat(os.Args[len(os.Args)-1])
	if err != nil || !os.SameFile(actualInfo, wantedInfo) {
		os.Exit(1)
	}
	fmt.Print("directory-ok")
	os.Exit(0)
}
