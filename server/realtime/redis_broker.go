package realtime

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisBroker struct {
	client       *redis.Client
	pubsub       *redis.PubSub
	incomingChan <-chan redis.Message
}

func NewRedisBroker(client *redis.Client) *RedisBroker {
	return &RedisBroker{
		client:       client,
		pubsub:       nil,
		incomingChan: make(<-chan redis.Message, 100),
	}
}

func (rb *RedisBroker) Initialize() error {
	return nil
}

// TODO: Improve this
func (rb *RedisBroker) Publish(channel string, message Envelope) error {
	ctx := context.Background()
	rb.client.Publish(ctx, channel, message)

	return nil
}

// TODO: Improve this
func (rb *RedisBroker) Subscribe(channel string) {
	ctx := context.Background()
	rb.client.Subscribe(ctx, channel)
}

// TODO: Improve this
func (rb *RedisBroker) Unsubscribe(channel string) {
	// ctx := context.Background()

}

// TODO: Improve this
func (rb *RedisBroker) ListenToSubscriptions() <-chan any {
	outChan := make(chan any, 100)
	go func() {
		for incomingMessage := range rb.incomingChan {
			outChan <- incomingMessage
		}
	}()

	return outChan
}
