package channels

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/google/uuid"
)

func New(ctx context.Context, name string, hubChannel chan realtime.Envelope) *Channel {
	ctx, cancel := context.WithCancel(ctx)
	return &Channel{
		Name:        name,
		Stream:      make(chan ChannelMessage, 100),
		Subscribers: make(map[uuid.UUID]bool),
		ctx:         ctx,
		cancel:      cancel,
		hubChannel:  hubChannel,
	}
}

func (c *Channel) Subscribe(uid uuid.UUID) bool {
	c.mu.Lock()
	c.Subscribers[uid] = true
	c.mu.Unlock()

	return true
}

func (c *Channel) Unsubscribe(uid uuid.UUID) bool {
	_, found := c.Subscribers[uid]
	if !found {
		return false
	}

	c.mu.Lock()
	delete(c.Subscribers, uid)
	c.mu.Unlock()

	return true
}

func (c *Channel) Get() ChannelMessage {
	return <-c.Stream
}

func (c *Channel) Put(message ChannelMessage) {
	c.Stream <- message
}

// TODO: Improve this
func (c *Channel) Process() {
	ticker := time.NewTicker(c.ProcessInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				fmt.Println("Channel process terminated via context completion")
				return
			case <-ticker.C:
				message := c.Get()
				rawBytes, err := json.Marshal(message)
				if err != nil {
					fmt.Println("Failed to marshall Channel Message to raw bytes")
					continue
				}

				for range c.Subscribers {
					envelope := realtime.Envelope{
						Header: realtime.Header{
							SenderID: message.UserID,
						},
						Data: rawBytes,
					}

					c.hubChannel <- envelope
				}
			}
		}
	}()
}

func (c *Channel) Stop() {
	c.cancel()
}
