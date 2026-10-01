package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	dialer := websocket.Dialer{}
	conn, response, err := dialer.DialContext(ctx, "ws://localhost:8000/ws", http.Header{})
	if err != nil {
		log.Fatalf("Failed to connect to websocket: %v", err)
	}
	defer response.Body.Close()

	res, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Println("Response:\n", string(res))

	// 1. Set up your deadlines and handlers exactly as you did
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	conn.SetPingHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		// Send the Pong back to the server
		return conn.WriteMessage(websocket.PongMessage, []byte(appData))
	})

	// Monitor context cancellation in a clean, separate goroutine
	go func() {
		<-ctx.Done()
		fmt.Println("Context cancelled, closing connection...")
		conn.Close() // This will force ReadMessage() below to unblock with an error
	}()

	// 2. Simple, non-blocking, clean read loop
	for {
		messageType, raw, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Connection closed or failed to read message: %v\n", err)
			// CRITICAL: Stop the loop immediately to avoid the repeated read panic
			break
		}

		// Refresh the read deadline when a normal data message is received too!
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))

		// Handle normal text/binary messages
		if messageType == websocket.TextMessage {
			fmt.Println(string(raw))
		}
	}
}
