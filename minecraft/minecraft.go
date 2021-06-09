package minecraft

import (
	"io"

	"../process"
)

// Handler handles minecraft server process
type Handler struct {
	process process.Process
}

// NewHandler makes new minecraft handler
func NewHandler() *Handler {
	return &Handler{process: process.Process{}}
}

// Start starts minectaft server
func (m *Handler) Start(javaPath string, opt ...string) error {
	opt = append(opt, "-jar", "server.jar", "nogui")
	if javaPath == "" {
		javaPath = "java"
	}
	return m.process.Start(javaPath, opt...)
}

// Interrupt send keyboard interrupt to the minecraft server
func (m *Handler) Interrupt() {
	m.process.Interrupt()
	return
}

// Kill send kill the minecraft server immediately
func (m *Handler) Kill() {
	m.process.Kill()
	return
}

// Pipes returns minecrafts process stdin / stdout / stderr pipes
//    these pipes are automatically closed when the process killed by Handler.Stop()
func (m *Handler) Pipes() (io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	return m.process.Pipes()
}
