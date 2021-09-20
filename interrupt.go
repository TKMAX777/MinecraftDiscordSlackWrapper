package main

import (
	"os"
	"os/signal"
	"syscall"
)

// setupCloseHandler setup actions when input keyboard interrupt
func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	Minecraft.Interrupt()
	os.Exit(0)
}
