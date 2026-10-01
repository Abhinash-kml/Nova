package realtime

import (
	"encoding/json"
	"time"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Uid    uuid.UUID
	send   chan Envelope
	conn   *websocket.Conn
	hub    *Hub
	pm     *PresenceManager
	config *config.WebsocketConfig
}

func NewClient(config *config.WebsocketConfig, uid uuid.UUID, connection *websocket.Conn, pm *PresenceManager, hub *Hub) *Client {
	return &Client{
		Uid:    uid,
		conn:   connection,
		pm:     pm,
		hub:    hub,
		send:   make(chan Envelope, 1000),
		config: config,
	}
}

// TODO: Implement this
func (c *Client) ReadIncoming() {
	defer func() {
		c.conn.Close()
		c.hub.Unregister(c)
	}()

	c.conn.SetReadLimit(c.config.MessageSize)
	c.conn.SetReadDeadline(time.Now().Add(c.config.PongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.config.PongWait))
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

			c.hub.Logger().Error("Reading message error", zap.Error(err))
			break
		}

		c.conn.SetReadDeadline(time.Now().Add(c.config.PongWait))

		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}

		// TODO: Maybe implement a pipeline for this ?

		// Parse incoming message
		var incomingEnvelope Envelope
		err = json.Unmarshal(payload, &incomingEnvelope)
		if err != nil {
			c.hub.logger.Error("Failed to parse incoming envelope", zap.Error(err))
			continue
		}

		// Drop incoming message if it exceeded its ttl (message can be delayed dudee to network issues)
		if time.Since(incomingEnvelope.Header.CreatedAt) >= incomingEnvelope.Header.TTL {
			continue
		}

		// If message type is Presence event - simply send it to Presence manager
		if incomingEnvelope.Header.Type == MessagePresence {
			c.pm.SetStatus(c.Uid, incomingEnvelope)
		}

		c.hub.Send(incomingEnvelope)
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
		case messageEnvelope, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteWait))

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.hub.Logger().Error("Failed to get writer for writing websocket message. Skipped some messages...", zap.Error(err))
				break Loop
			}

			encoder := json.NewEncoder(writer)
			if encoder == nil {
				c.hub.Logger().Error("Failed to create json encoder. Websocker messages skipped")
				continue
			}
			err = encoder.Encode(messageEnvelope)
			if err != nil {
				c.hub.Logger().Error("Failed to encode envelope type to json. Skipped message", zap.Error(err))
				continue
			}

			// Batch all messages in the send channel using newline \n character
			len := len(c.send)
			for range len {
				writer.Write([]byte{'\n'})
				message := <-c.send
				encoder := json.NewEncoder(writer)
				err := encoder.Encode(message)
				if err != nil {
					c.hub.Logger().Error("Failed to encode envelope type to json. Skipped message", zap.Error(err))
					continue
				}
			}

			if err := writer.Close(); err != nil {
				c.hub.Logger().Error("Failed to write message using websocket writer", zap.Error(err))
				break Loop
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				c.hub.Logger().Error("Failed to send ping message to client", zap.Error(err))
				return
			}
		}
	}
}

func (c *Client) Send(message Envelope) {
	c.send <- message
}
