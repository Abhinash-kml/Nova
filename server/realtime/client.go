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
	send chan []byte
	conn *websocket.Conn
	hub  *Hub
}

func NewClient(uid uuid.UUID, connection *websocket.Conn, hub *Hub) *Client {
	return &Client{
		Uid:  uid,
		conn: connection,
		hub:  hub,
		send: make(chan []byte, 1000),
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
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway) {
				return
			}
		}

		c.hub.Send(raw)
	}
}

// TODO: Implement this
func (c *Client) ProcessOutgoing() {
	// Write loop
}
