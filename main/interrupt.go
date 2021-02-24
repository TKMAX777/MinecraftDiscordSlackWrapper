package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// setupCloseHandler setup actions when input keyboard interrupt
func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	errorHandle(fmt.Errorf("[Server] KeyboardInterrupt"))
	errorHandle(fmt.Errorf("Now killing minecraft server"))
	Minecraft.Interrupt()
	errorHandle(fmt.Errorf("[Server] Finished! GoodBye"))
	os.Exit(0)
}
