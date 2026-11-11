package messagebroker

import (
	"context"
	"sync"
)

type InMemoryMessageBroker struct {
	mu       sync.RWMutex
	channels map[string]chan map[string]any
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewInMemoryMessageBroker(ctx context.Context, size int) *InMemoryMessageBroker {
	innerCtx, cancelFunc := context.WithCancel(ctx)
	return &InMemoryMessageBroker{
		channels: make(map[string]chan map[string]any, size),
		ctx:      innerCtx,
		cancel:   cancelFunc,
	}
}

func (em *InMemoryMessageBroker) Initialize() {}

// Helper method to safely ensure a channel exists under a write lock
func (em *InMemoryMessageBroker) getOrCreateChannel(channel string) chan map[string]any {
	// 1. First check with a lightweight Read lock
	em.mu.RLock()
	ch, found := em.channels[channel]
	em.mu.RUnlock()
	if found {
		return ch
	}

	// 2. Double check under a full Write lock to prevent race conditions
	em.mu.Lock()
	defer em.mu.Unlock()

	// Double check constraint in case another thread created it while we upgraded locks
	if ch, found = em.channels[channel]; !found {
		// Using a larger buffer size prevents publishers from blocking immediately
		ch = make(chan map[string]any, 500)
		em.channels[channel] = ch
	}
	return ch
}

func (em *InMemoryMessageBroker) Produce(channel string, payload map[string]any) {
	ch := em.getOrCreateChannel(channel)

	// Safely emit without locking the entire map during network/channel blockage
	select {
	case ch <- payload:
	case <-em.ctx.Done():
		return
	}
}

func (em *InMemoryMessageBroker) Consume(channel string) <-chan map[string]any {
	ch := em.getOrCreateChannel(channel)
	outCh := make(chan map[string]any, 100)

	go func() {
		defer close(outCh)
		for {
			select {
			case <-em.ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return // Underlying channel was closed
				}

				// Safely forward down the stream pipeline
				select {
				case outCh <- msg:
				case <-em.ctx.Done():
					return
				}
			}
		}
	}()

	return outCh
}
