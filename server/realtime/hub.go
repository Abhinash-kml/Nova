package realtime

import "fmt"

type Hub struct {
	// store SessionStore
	register   chan *Client
	unregister chan *Client
	send       chan Envelope
	incoming   chan Envelope
	broadcast  chan Envelope

	// pubsub PubSub
	// pool WorkerPool
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client, 100),
		unregister: make(chan *Client, 100),
		send:       make(chan Envelope, 100),
		incoming:   make(chan Envelope, 100),
		broadcast:  make(chan Envelope, 100),
	}
}

func (h *Hub) Initialize() {

}

func (h *Hub) Run() {
	fmt.Println("Starting realtime hub")

	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)
		case client := <-h.unregister:
			h.handleUnregister(client)
		case message := <-h.send:
			h.handleSend(message)
		case message := <-h.broadcast:
			h.handleBroadcast(message)
		case message := <-h.incoming:
			h.handleIncoming(message)
		}
	}
}

func (h *Hub) Send(message Envelope) {
	h.send <- message
}

func (h *Hub) handleSend(message Envelope) {

}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) handleRegister(client *Client) {

}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) handleUnregister(client *Client) {

}

func (h *Hub) Broadcast(message Envelope) {
	h.broadcast <- message
}

func (h *Hub) handleBroadcast(message Envelope) {

}

func (h *Hub) handleIncoming(message Envelope) {

}

func (h *Hub) enrichMessage(message Envelope) {

}
