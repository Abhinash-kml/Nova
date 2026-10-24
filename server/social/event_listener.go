package social

import (
	"context"
	"encoding/json"
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

			var cMessage realtime.ChatMessage
			err := json.Unmarshal(message.Data, &cMessage)
			if err != nil {
				ml.logger.Error("Failed to unmarshall chat message", zap.Error(err))
				continue
			}

			err = ml.service.CreateConversationSendMessage(context.Background(), cMessage, message.Header.SenderID, message.Header.ReceiverID)
			if err != nil {
				ml.logger.Error("Failed to persist merssage in db", zap.Error(err))
			}
		}
	}()
}
