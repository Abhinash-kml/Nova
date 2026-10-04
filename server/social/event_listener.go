package social

import (
	"context"
	"fmt"

	"github.com/abhinash-kml/nova/server/realtime"
)

type MessageListener struct {
	service Service
	queue   chan realtime.Envelope
}

func NewMessageListener(service Service, queue chan realtime.Envelope) *MessageListener {
	return &MessageListener{
		service: service,
		queue:   queue,
	}
}

func (ml *MessageListener) Listen() {
	go func() {
		for {
			message, ok := <-ml.queue
			if !ok {
				return
			}

			fmt.Println("Got message:", message.Header.SenderID)

			ml.service.SendMessage(context.Background(), message)
		}
	}()
}
