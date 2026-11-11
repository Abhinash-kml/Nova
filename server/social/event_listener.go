package social

import (
	"context"
	"time"

	"github.com/abhinash-kml/nova/server/messagebroker"
	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MessageListener struct {
	ctx     context.Context
	cancel  context.CancelFunc
	service Service
	logger  *zap.Logger
}

func NewMessageListener(ctx context.Context, service Service, l *zap.Logger) *MessageListener {
	innerCtx, cancel := context.WithCancel(ctx)
	return &MessageListener{
		logger:  l,
		service: service,
		ctx:     innerCtx,
		cancel:  cancel,
	}
}

func (ml *MessageListener) Listen() {
	go func() {
		messageBroker := messagebroker.GetGlobal()
		incomingChan := messageBroker.Consume("stream:chatmessages")
		for {
			select {
			case message := <-incomingChan:
				conversationType := message["conversation_type"].(int)
				conversationID := message["chat_id"].(string)
				messageID := message["message_id"].(string)
				senderID := message["user_id"].(string)
				targetID := message["target_id"].(string)
				payload := message["payload"].(string)
				createdAt := message["created_at"].(time.Time)
				status := message["status"].(realtime.MessageStatus)
				deleted := message["deleted"].(bool)

				zap.L().Sugar().Debugf(
					"ConverationType: %d\nConversationID: %s\nMessageID: %s\nSenderID: %s\nTargetID: %s\nPayload: %s\nStatus: %s\nDeleted: %s",
					conversationType,
					conversationID,
					messageID,
					senderID,
					targetID,
					payload,
					status,
					deleted,
				)

				cMessage := realtime.ChatMessage{
					ConversationType: conversationType,
					ChatId:           uuid.MustParse(conversationID),
					MessageId:        uuid.MustParse(messageID),
					Body:             payload,
					Deleted:          deleted,
					Status:           status,
					CreatedAt:        createdAt,
				}

				err := ml.service.CreateConversationSendMessage(
					context.Background(),
					cMessage,
					uuid.MustParse(senderID),
					uuid.MustParse(targetID),
				)
				if err != nil {
					ml.logger.Error("Failed to persist merssage in db", zap.Error(err))
				}
			case <-ml.ctx.Done():
				zap.L().Sugar().Debug("Social event listener stopped due to context cancellation")
				return
			}
		}
	}()
}
