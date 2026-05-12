package realtime

type Hub struct {
	// store SessionStore
	register   chan *Client
	unregister chan *Client
	send       chan []byte
	incoming   chan []byte
	broadcast  chan []byte

	// pubsub PubSub
}

func (h *Hub) Initialize() {

}

func (h *Hub) Run() {

}

func (h *Hub) Send() {

}

func (h *Hub) handleSend() {

}

func (h *Hub) Register() {

}

func (h *Hub) handleRegister() {

}

func (h *Hub) Unregister(client *Client) {

}

func (h *Hub) handleUnregister() {

}

func (h *Hub) Broadcast() {

}

func (h *Hub) handleBroadcast() {

}

func (h *Hub) handleIncoming() {

}

func (h *Hub) enrichMessage() {

}
