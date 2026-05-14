package realtime

import "fmt"

type Hub struct {
	// store SessionStore
	register   chan *Client
	unregister chan *Client
	send       chan Envelope
	incoming   chan Envelope
	broadcast  chan Envelope

	broker   RealtimeBroker
	registry SessionStore
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
	fmt.Println("Renning realtime hub with 5 gorotines")

	for range 5 {
		go func() {
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
		}()
	}
}

// Send the message to the send channel
func (h *Hub) Send(message Envelope) {
	h.send <- message
}

// 1. Extract the receiver id from envelope header
// 2. Check if user exists locally on this node
// 2.1. If yes send to its send channel, client goroutine will handle writing to socket
// 2.2. If no send to the pubsub using reciever id
func (h *Hub) handleSend(message Envelope) {
	receiverId := message.Header.ReceiverID
	if h.registry.Exists(receiverId) {
		client := h.registry.Get(receiverId)
		client.Send(message)
		return
	}

	err := h.broker.Publish(receiverId.String(), message)
	if err != nil {
		// Handle
	}
}

// 1. Send Client info to the register channel
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// 1. Add client to the session store
// 2. Subscribe to the client's userid on pubsub to recieve incoming messsage for client
func (h *Hub) handleRegister(client *Client) {
	h.registry.Add(client)
	h.broker.Subscribe(client.Uid.String())
}

// 1. Send Client info to the unregister channel
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// 1. Remove client from the session store
// 2. unsubscribe from client's userid in pubsub to stop receiving message for the client
func (h *Hub) handleUnregister(client *Client) {
	h.registry.Remove(client)
	h.broker.Unsubscribe(client.Uid.String())
}

// 1. Send envelope to the broadcast channel
func (h *Hub) Broadcast(message Envelope) {
	h.broadcast <- message
}

// 1. Loop through all the users in the session store and send then the envelope
func (h *Hub) handleBroadcast(message Envelope) {
	h.registry.ForEach(func(c *Client) {
		c.Send(message)
	})
}

// 1. Send the incoming message to the send channel
func (h *Hub) handleIncoming(message Envelope) {

}

func (h *Hub) enrichMessage(message Envelope) {

}

func (h *Hub) SendChannel() chan Envelope {
	return h.send
}
