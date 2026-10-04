package social

import (
	"context"
	"fmt"

	"github.com/abhinash-kml/nova/server/realtime"
	"go.uber.org/zap"
)

type MessageListener struct {
	service Service
	queue   chan realtime.Envelope
	logger  *zap.Logger
}

func NewMessageListener(service Service, queue chan realtime.Envelope, l *zap.Logger) *MessageListener {
	return &MessageListener{
		logger:  l,
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

			err := ml.service.SendMessage(context.Background(), message)
			if err != nil {
				ml.logger.Error("Failed to persist merssage in db", zap.Error(err))
			}
		}
	}()
}
