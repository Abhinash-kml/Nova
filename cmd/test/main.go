package main

import (
	"os"
	"os/signal"
)

func main() {
	// Operations
	// PerformOAuthGoogleLibrary()
	PerformOAuthGoogleCustom()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
