package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/events"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Uid            uuid.UUID
	queuedMessages chan Envelope
	conn           *websocket.Conn
	config         *config.WebsocketConfig
}

func NewClient(config *config.WebsocketConfig, uid uuid.UUID, connection *websocket.Conn) *Client {
	return &Client{
		Uid:            uid,
		conn:           connection,
		queuedMessages: make(chan Envelope, 1000),
		config:         config,
	}
}

func (c *Client) ReadIncoming() {
	defer func() {
		c.conn.Close()
		GetGlobalHub().Unregister(c)
	}()

	c.conn.SetReadLimit(c.config.MessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(c.config.PingInterval + c.config.PongWait))
	if err != nil {
		zap.L().Error("Failed to set read deadline", zap.Error(err))
		return
	}

	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.config.PingInterval + c.config.PongWait))
		GetPresenceManager().RefreshSession(context.Background(), c.Uid.String())

		zap.L().Debug("Pong received from client", zap.String("uid", c.Uid.String()))
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

		c.conn.SetReadDeadline(time.Now().Add(c.config.PingInterval + c.config.PongWait))

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

		zap.L().Sugar().Debug("========== BEGIN MESSAGE ==========")

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "   ")
		encoder.Encode(incomingEnvelope)

		c.ProcessMessage(incomingEnvelope)

		zap.L().Sugar().Debug("========== END MESSAGE ==========")
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

			err := c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteWait))
			if err != nil {
				zap.L().Error("Failed to set write deadline", zap.Error(err))
				return
			}

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
			err := c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteWait))
			if err != nil {
				zap.L().Error("Failed to set write deadline", zap.Error(err))
				return
			}
			err = c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				zap.L().Error("Failed to send ping message to client", zap.Error(err))
				return
			}

			// err = c.conn.SetWriteDeadline(time.Time{}) // Expirement
			// if err != nil {
			// 	zap.L().Error("Failed to set write deadline", zap.Error(err))
			// 	return
			// }

			zap.L().Debug("Sent ping to client", zap.String("uid", c.Uid.String()))
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
		zap.L().Sugar().Debug("Sending receipt message to hub")

		GetGlobalHub().Send(message)
		return
	}

	// 2. Persist
	err = c.PersistMessage(message)
	if err != nil {
		zap.L().Error("Failed to persist message", zap.Error(err))
		return
	}

	zap.L().Sugar().Debug("After Persisting")
	zap.L().Sugar().Debug("Sending acknowledgement")

	// 3. Send Acknowledgement receipt
	err = c.SendAcknowledgement(message)
	if err != nil {
		zap.L().Error("Failed to send acknowledgement", zap.Error(err))
	}

	// 4. Forward it to hub for realtime forwarding
	GetGlobalHub().Send(message)
	zap.L().Sugar().Debug("After fowarding to hub")
}

func (c *Client) PersistMessage(message Envelope) error {
	var cMessage ChatMessage
	err := json.Unmarshal(message.Data, &cMessage)
	if err != nil {
		return fmt.Errorf("unmarshalling chat message: %w", err)
	}

	eventsManager := events.GetGlobalEventsManager()
	eventsManager.Produce("stream:chatmessages", map[string]any{
		"conversation_type": cMessage.ConversationType,
		"chat_id":           cMessage.ChatId.String(),
		"message_id":        cMessage.MessageId.String(),
		"user_id":           message.Header.SenderID.String(),
		"target_id":         message.Header.ReceiverID.String(),
		"payload":           cMessage.Body,
		"deleted":           cMessage.Deleted,
		"status":            cMessage.Status,
		"created_at":        message.Header.CreatedAt,
	})
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
