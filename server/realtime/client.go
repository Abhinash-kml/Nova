package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Uid            uuid.UUID
	queuedMessages chan Envelope
	conn           *websocket.Conn
	config         *config.WebsocketConfig
	persister      MessagePersister
}

func NewClient(config *config.WebsocketConfig, uid uuid.UUID, connection *websocket.Conn) *Client {
	return &Client{
		Uid:            uid,
		conn:           connection,
		queuedMessages: make(chan Envelope, 1000),
		config:         config,
		persister:      GetLocalMessagePersistor(),
	}
}

// TODO: Implement this
func (c *Client) ReadIncoming() {
	defer func() {
		c.conn.Close()
		GetGlobalHub().Unregister(c)
	}()

	c.conn.SetReadLimit(c.config.MessageSize)
	c.conn.SetReadDeadline(time.Now().Add(c.config.PongWait))

	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.config.PongWait))
		GetPresenceManager().RefreshSession(context.Background(), c.Uid.String())
		return nil
	})

	// Read loop
	for {
		messageType, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway) {
				return
			}

			zap.L().Error("Reading message error", zap.Error(err))
			break
		}

		c.conn.SetReadDeadline(time.Now().Add(c.config.PongWait))

		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}

		// Parse incoming message
		var incomingEnvelope Envelope
		err = json.Unmarshal(payload, &incomingEnvelope)
		if err != nil {
			zap.L().Error("Failed to parse incoming envelope", zap.Error(err))
			continue
		}

		fmt.Println("========== BEGIN MESSAGE ==========")

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "   ")
		encoder.Encode(incomingEnvelope)

		c.ProcessMessage(incomingEnvelope)

		fmt.Println("========== END MESSAGE ==========")
	}
}

// TODO: Subjected to improvement
func (c *Client) ProcessOutgoing() {
	// Write loop
	ticker := time.NewTicker(c.config.PingInterval)
	defer ticker.Stop()

Loop:
	for {
		select {
		case messageEnvelope, ok := <-c.queuedMessages:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteWait))

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				zap.L().Error("Failed to get writer for writing websocket message. Skipped some messages...", zap.Error(err))
				break Loop
			}

			encoder := json.NewEncoder(writer)
			if encoder == nil {
				zap.L().Error("Failed to create json encoder. Skipped message")
				continue
			}

			err = encoder.Encode(messageEnvelope)
			if err != nil {
				zap.L().Error("Failed to encode outgoing envelope. Skipped message", zap.Error(err))
				continue
			}

			// Batch all messages in the send channel using newline \n character
			len := len(c.queuedMessages)
			for range len {
				writer.Write([]byte{'\n'})
				message := <-c.queuedMessages
				encoder := json.NewEncoder(writer)
				err := encoder.Encode(message)
				if err != nil {
					zap.L().Error("Failed to encode envelope type to json. Skipped message", zap.Error(err))
					continue
				}
			}

			if err := writer.Close(); err != nil {
				zap.L().Error("Failed to write message using writer", zap.Error(err))
				break Loop
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				zap.L().Error("Failed to send ping message to client", zap.Error(err))
				return
			}
		}
	}
}

func (c *Client) Send(message Envelope) {
	c.queuedMessages <- message
}

func (c *Client) ProcessMessage(message Envelope) {
	var err error

	// 1. Filter using meta-data
	if time.Since(message.Header.CreatedAt) > message.Header.TTL {
		return
	}

	// If message type is Presence event - simply send it to Presence manager
	if message.Header.Type == MessagePresenceAction {
		GetPresenceManager().Process(message)
		return
	}

	// If message type is lobby event - send it to lobby manager
	if message.Header.Type == MessageLobbyEvent {
		// GetLobbyManager().Send(message)
	}

	// If message type is Chat Receipt - dont persist just forward
	if message.Header.Type == MessageReceipt {
		fmt.Println("Sending receipt message to hub")
		GetGlobalHub().Send(message)
		return
	}

	// 2. Persist
	err = c.PersistMessage(message)
	if err != nil {
		zap.L().Error("Failed to persist message", zap.Error(err))
		return
	}

	fmt.Println("After Persisting")

	// 3. Send Acknowledgement receipt
	err = c.SendAcknowledgement(message)
	if err != nil {
		zap.L().Error("Failed to send acknowledgement", zap.Error(err))
	}
	// fmt.Println("Sending acknowledgemenet")

	// 4. Forward it to hub for realtime forwarding
	GetGlobalHub().Send(message)
	fmt.Println("After fowarding to hub")
}

func (c *Client) PersistMessage(message Envelope) error {
	err := c.persister.Persist(message)
	if err != nil {
		return fmt.Errorf("persisting message: %w", err)
	}

	return nil
}

func (c *Client) SendAcknowledgement(message Envelope) error {
	var data ChatMessage
	json.Unmarshal(message.Data, &data)

	receipt := ChatReceipt{
		MessageId: data.MessageId,
		Status:    StatusSent,
	}

	raw, err := json.Marshal(receipt)
	if err != nil {
		return fmt.Errorf("marshalling reciept: %w", err)
	}

	envelope := Envelope{
		Header: Header{
			Type:       MessageReceipt,
			SourceID:   uuid.New(),
			ReceiverID: message.Header.SenderID,
		},
		Data: json.RawMessage(raw),
	}

	c.Send(envelope)

	return nil
}
