package realtime

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	WriteWait    = 10 * time.Second
	PongWait     = 5 * time.Second
	PingInterval = 5 * time.Second
	MaxMessages  = 1024
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

	c.conn.SetReadLimit(int64(c.config.MessageSize))
	c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.config.PongWait)))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.config.PongWait)))
		return nil
	})

	// Read loop
	for {
		var incomingEnvelope Envelope
		err := c.conn.ReadJSON(&incomingEnvelope)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway) {
				return
			}
		}

		// TODO: Maybe implement a pipeline for this ?

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
	ticker := time.NewTicker(PingInterval)
	defer ticker.Stop()

	for {
		select {
		case messageEnvelope, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(time.Duration(c.config.WriteWait)))

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println("Failed to get writer for writing websocket message. Skipped some messages...")
				break
			}

			encoder := json.NewEncoder(writer)
			if encoder == nil {
				fmt.Println("Failed to create json encoder. Websocker messages skipped")
				continue
			}
			err = encoder.Encode(messageEnvelope)
			if err != nil {
				fmt.Println("Failed to encode envelope type to json. Skipped message")
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
					fmt.Println("Failed to encode envelope type to json. Skipped message")
					continue
				}
			}

			if err := writer.Close(); err != nil {
				fmt.Println("Failed to write message using websocket writer")
				continue
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(time.Duration(c.config.WriteWait)))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				// Handle error here
				// Log
				return
			}
		}

	}
}

func (c *Client) Send(message Envelope) {
	c.send <- message
}
