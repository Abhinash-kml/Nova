package realtime

import (
	"time"

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
	Uid  uuid.UUID
	send chan Envelope
	conn *websocket.Conn
	hub  *Hub
}

func NewClient(uid uuid.UUID, connection *websocket.Conn, hub *Hub) *Client {
	return &Client{
		Uid:  uid,
		conn: connection,
		hub:  hub,
		send: make(chan Envelope, 1000),
	}
}

// TODO: Implement this
func (c *Client) ReadIncoming() {
	defer func() {
		c.conn.Close()
		c.hub.Unregister(c)
	}()

	c.conn.SetReadLimit(1024)
	c.conn.SetReadDeadline(time.Now().Add(PongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(PongWait))
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

			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			err := c.conn.WriteJSON(messageEnvelope)
			if err != nil {
				// Handle
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				// Handle error here
				// Log
				return
			}
		}

	}
}
