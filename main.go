package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/PauSabatesC/congo/cmd"
)

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
	}()
}

func main() {
	cmd.Execute()
}
