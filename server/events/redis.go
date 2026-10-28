package events

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisEventsManager struct {
	redis  *redis.Client
	ctx    context.Context
	cancal context.CancelFunc
}

func NewRedisEventsManager(ctx context.Context, redis *redis.Client) *RedisEventsManager {
	innerCtx, cancelFunc := context.WithCancel(ctx)
	return &RedisEventsManager{
		redis:  redis,
		ctx:    innerCtx,
		cancal: cancelFunc,
	}
}

func (re *RedisEventsManager) Initialize() {

}

func (re *RedisEventsManager) Produce(channel string, payload map[string]any) {
	fmt.Println("Got event. producing....")
	err := re.redis.XAdd(re.ctx, &redis.XAddArgs{
		Stream: channel,
		ID:     "*",
		Values: payload,
	}).Err()
	if err != nil {
		fmt.Println("Failed to send redis stream event. Error:", err)
	}
}

func (re *RedisEventsManager) Consume(channel string) <-chan map[string]any {
	outCh := make(chan map[string]any, 100)

	go func() {
		lastID := "$"
		for {
			select {
			case <-re.ctx.Done():
				fmt.Println("Stopped consuming from", channel, "channel due to context cancellation")
				return
			default:
			}

			streams, err := re.redis.XRead(re.ctx, &redis.XReadArgs{
				Streams: []string{channel, lastID},
				Block:   0,
			}).Result()
			if err != nil {
				fmt.Println("Failed to consume from stream:", channel, "Error:", err)
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					lastID = message.ID
					select {
					case outCh <- message.Values:
					case <-re.ctx.Done():
						return
					}
				}
			}
		}
	}()

	return outCh
}
