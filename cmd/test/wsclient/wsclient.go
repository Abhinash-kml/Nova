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
	"slices"
	"strings"
	"time"

	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	SERVER_PING_INTERVAL = time.Second * 15
	READ_DEADLINE        = SERVER_PING_INTERVAL + (time.Second * 5)
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
	conn.SetReadDeadline(time.Now().Add(READ_DEADLINE))

	conn.SetPingHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(READ_DEADLINE))
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
		conn.SetReadDeadline(time.Now().Add(READ_DEADLINE))

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
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte(name))
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
	fmt.Println("What you wanna do ?")
	fmt.Println("1. Messaging [ 1-1 / Group ] ")
	fmt.Println("2. Operations [ sub, unsub, status ]")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	choice := scanner.Text()
	switch choice {
	case "1":
		PerformMessaging(conn, senderID)
	case "2":
		PerformOperations(conn, senderID)
	}
}

func PerformMessaging(conn *websocket.Conn, senderID uuid.UUID) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		chatType := parts[0]
		targetId := NameToUUID(parts[1])
		body := parts[2]

		// 1:username_name: message
		// 2:group_name: message

		envelope := ConstuctChatMessage(senderID, targetId, body, chatType)

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

func PerformOperations(conn *websocket.Conn, senderID uuid.UUID) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		operation := parts[0]
		cleaned := strings.TrimSpace(parts[1])
		targets := strings.Split(cleaned, " ")

		envelope := ConstructOperationsMessage(senderID, operation, targets)

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

func ConstuctChatMessage(senderID, receiverID uuid.UUID, message string, chatType string) realtime.Envelope {
	var ConversationID uuid.UUID

	// If its a direct chat create conversation id by combining sender and receiver id and geenrate uuidv5
	if chatType == "1" {
		temp := []string{senderID.String(), receiverID.String()}
		slices.Sort(temp)
		combined := temp[0] + ":" + temp[1]
		ConversationID = uuid.NewSHA1(uuid.NameSpaceDNS, []byte(combined))
	} else { // For group chat the conversation id is the reciver id generated by creating uuidv5 from group name
		ConversationID = receiverID
	}

	fmt.Println("Conversation ID:", ConversationID.String())

	chatMessage := realtime.ChatMessage{
		ChatId:    ConversationID,
		MessageId: uuid.New(),
		Body:      message,
		CreatedAt: time.Now(),
	}
	if chatType == "1" {
		chatMessage.ConversationType = 1
	} else {
		chatMessage.ConversationType = 2
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

func ConstructOperationsMessage(actionTakerID uuid.UUID, action string, targets []string) realtime.Envelope {
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
