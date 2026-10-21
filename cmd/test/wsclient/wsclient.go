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

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	userID := SetupUser()

	dialer := websocket.Dialer{}
	header := http.Header{}
	header.Add("userid", userID.String())
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

	go ReadFromStdIn(conn, userID)

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
			var envelope realtime.Envelope
			json.Unmarshal(raw, &envelope)

			buffer := bytes.NewBuffer(make([]byte, len(raw)))
			json.Indent(buffer, raw, " ", "   ")

			fmt.Println(buffer)

			// Send back acknowledgement
			if envelope.Header.Type == realtime.MessageChat {
				var cd realtime.ChatMessage
				json.Unmarshal(envelope.Data, &cd)

				SendAcknowledgement(conn, cd.ChatId, cd.MessageId, userID, envelope.Header.SenderID)
			}
		}
	}
}

func SetupUser() uuid.UUID {
	fmt.Println("Enter username:")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()
	return NameToUUID(username)
}

func NameToUUID(name string) uuid.UUID {
	namespace := uuid.NameSpaceDNS
	return uuid.NewSHA1(namespace, []byte(name))
}

func SendAcknowledgement(conn *websocket.Conn, chatID, messageID, senderID, receiverID uuid.UUID) {
	header := realtime.Header{
		Type:       realtime.MessageReceipt,
		SenderID:   senderID,
		ReceiverID: receiverID,
		TTL:        time.Second * 5,
	}

	data := realtime.ChatReceipt{
		ChatId:    chatID,
		MessageId: messageID,
		Status:    realtime.StatusRead,
	}

	bytes, _ := json.Marshal(data)

	envelope := realtime.Envelope{
		Header: header,
		Data:   bytes,
	}

	finalBytes, _ := json.Marshal(envelope)

	err := conn.WriteMessage(websocket.TextMessage, finalBytes)
	if err != nil {
		fmt.Println("Failed to send acknowledgement message:", err)
	}

	fmt.Println("Acknowledgement sent")
}

func ReadFromStdIn(conn *websocket.Conn, senderID uuid.UUID) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var envelope realtime.Envelope
		line := scanner.Text()
		parts := strings.Split(line, " ")

		// Chat message
		if parts[0] != "sub" && parts[0] != "unsub" && parts[0] != "status" {
			reciverID := NameToUUID(parts[0])
			envelope = ConstuctChatMessage(senderID, reciverID, parts[1])
		} else {
			// Presence action
			envelope = ConstructPresenceSubUnsub(senderID, parts[0], parts[1:])
		}

		writer, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			fmt.Println("Failed to get next writer to write messages")
			return
		}
		defer writer.Close()

		err = json.NewEncoder(writer).Encode(envelope)
		if err != nil {
			fmt.Println("Failed to encode and send message", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read: %v", err)
	}
}

func ConstuctChatMessage(senderID, receiverID uuid.UUID, message string) realtime.Envelope {
	chatID, _ := uuid.Parse("912426f3-ccb0-5dcd-b937-8aebce8f2ecd")
	chatMessage := realtime.ChatMessage{
		ChatId:    chatID,
		MessageId: uuid.New(),
		Body:      message,
		CreatedAt: time.Now(),
	}

	header := realtime.Header{
		Type:       realtime.MessageChat,
		SenderID:   senderID,
		ReceiverID: receiverID,
		CreatedAt:  time.Now(),
		TTL:        time.Duration(time.Hour * 12),
	}

	raw, _ := json.Marshal(chatMessage)
	envelope := realtime.Envelope{
		Header: header,
		Data:   json.RawMessage(raw),
	}

	return envelope
}

func ConstructPresenceSubUnsub(actionTakerID uuid.UUID, action string, targets []string) realtime.Envelope {
	presenceActionEvent := realtime.PresenceAction{
		UserID:   actionTakerID,
		ActionAt: time.Now(),
	}

	switch action {
	case "sub":
		presenceActionEvent.Action = "sub"
	case "unsub":
		presenceActionEvent.Action = "unsub"
	case "status":
		presenceActionEvent.Action = "status"
	}

	for i := range targets {
		ID := NameToUUID(targets[i])
		presenceActionEvent.Targets = append(presenceActionEvent.Targets, ID)
	}

	marshalledAction, _ := json.Marshal(presenceActionEvent)

	envelope := realtime.Envelope{
		Header: realtime.Header{
			SenderID:  actionTakerID,
			Type:      realtime.MessagePresenceAction,
			CreatedAt: time.Now(),
			TTL:       time.Second * 30,
		},
		Data: marshalledAction,
	}

	return envelope
}
