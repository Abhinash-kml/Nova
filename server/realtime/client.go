package realtime

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

	// Read loop
}

// TODO: Implement this
func (c *Client) ProcessOutgoing() {
	// Write loop
}
