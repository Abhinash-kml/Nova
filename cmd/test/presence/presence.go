package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	SERVER_PING_INTERVAL time.Duration = time.Second * 15
	READ_DEADLINE                      = SERVER_PING_INTERVAL + (time.Second * 5)
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	connections := ConnectNClients(ctx)

	go ReadAndSendMessages(connections)

	<-ctx.Done()
}

func ConnectNClients(ctx context.Context) map[uuid.UUID]*websocket.Conn {
	users := SetupUsers()

	var connections map[uuid.UUID]*websocket.Conn
	connections = make(map[uuid.UUID]*websocket.Conn, len(users))

	for i := range users {
		go func() {
			dialer := websocket.Dialer{}
			header := http.Header{}
			header.Add("userid", users[i].String())
			conn, response, err := dialer.DialContext(ctx, "ws://localhost:8000/ws", header)
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
			conn.SetReadDeadline(time.Now().Add(READ_DEADLINE))

			conn.SetPingHandler(func(appData string) error {
				conn.SetReadDeadline(time.Now().Add(READ_DEADLINE))
				// Send the Pong back to the server
				return conn.WriteMessage(websocket.PongMessage, []byte(appData))
			})

			connections[users[i]] = conn

			ListenForMessages(conn)
		}()
	}

	// Monitor context cancellation in a clean, separate goroutine
	go func() {
		<-ctx.Done()
		fmt.Println("Context cancelled, closing connection all...")
		for i := range connections {
			connections[i].Close() // This will force ReadMessage() below to unblock with an error
		}
	}()

	return connections
}

func ListenForMessages(conn *websocket.Conn) {
	// 2. Simple, non-blocking, clean read loop
	for {
		messageType, raw, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Connection closed or failed to read message: %v\n", err)
			// CRITICAL: Stop the loop immediately to avoid the repeated read panic
			break
		}

		// Refresh the read deadline when a normal data message is received too!
		conn.SetReadDeadline(time.Now().Add(READ_DEADLINE))

		// Handle normal text/binary messages
		if messageType == websocket.TextMessage {
			var envelope realtime.Envelope
			json.Unmarshal(raw, &envelope)

			var cd realtime.ChatMessage
			json.Unmarshal(envelope.Data, &cd)

			buffer := bytes.NewBuffer(make([]byte, len(raw)))
			json.Indent(buffer, raw, " ", "   ")

			fmt.Println(buffer)

			if envelope.Header.Type == realtime.MessagePresence {
				fmt.Println("Presence message received")
			}

			// Send back acknowledgement
			// if envelope.Header.Type == realtime.MessageChat {
			// 	SendAcknowledgement(conn, cd.ChatId, cd.MessageId, userID, envelope.Header.SenderID)
			// }
		}
	}
}

func ReadAndSendMessages(connections map[uuid.UUID]*websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		parts := strings.Split(text, ":")
		userID := NameToUUID(parts[0])
		status := strings.TrimSpace(parts[1])

		// name: online
		// name: offline
		// name: playing in erangel

		presenceEvent := realtime.PresenceAction{
			UserID:   userID,
			Action:   "set",
			ActionAt: time.Now(),
		}

		switch {
		case status == "online":
			presenceEvent.Data = map[string]string{
				"status": "online",
			}
		case status == "offline":
			presenceEvent.Data = map[string]string{
				"status": "offline",
			}
		default:
			fmt.Println("Custom Status:", status)
			presenceEvent.Data = map[string]string{
				"status":        "custom",
				"custom_status": status,
			}
		}

		fmt.Printf("%s: %s\n", userID, status)

		bytes, _ := json.Marshal(presenceEvent)

		envelope := realtime.Envelope{
			Header: realtime.Header{
				SenderID:  userID,
				Type:      realtime.MessagePresenceAction,
				CreatedAt: time.Now(),
				TTL:       time.Second * 30,
			},
			Data: bytes,
		}

		marshalledEnvelope, _ := json.Marshal(envelope)

		err := connections[userID].WriteMessage(websocket.TextMessage, marshalledEnvelope)
		if err != nil {
			fmt.Printf("Failed to send message to %s[%s]: %v", parts[0], userID.String(), err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to scan message from stdin:", err)
	}
}

func SetupUsers() []uuid.UUID {
	fmt.Println("Enter usernames:")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	usernames := scanner.Text()
	parts := strings.Split(usernames, " ")

	var userIDs []uuid.UUID
	for i := range parts {
		userIDs = append(userIDs, NameToUUID(parts[i]))
	}
	return userIDs
}

func NameToUUID(name string) uuid.UUID {
	namespace := uuid.NameSpaceDNS
	return uuid.NewSHA1(namespace, []byte(name))
}
