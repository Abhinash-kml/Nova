package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Listen for interrupt & kill signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Global context for passing to all services
	globalCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Perform our task

	// Block untill our signal is trigerred
	<-signalChan

	// Gracefully shutdown all services by calling cancel() of global context
	cancel()
}
