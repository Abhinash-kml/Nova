package realtime

import (
	"sync"

	"go.uber.org/zap"
)

var (
	once      sync.Once
	persistor *LocalMessagePersister
)

func GetLocalMessagePersistor() *LocalMessagePersister {
	once.Do(func() {
		if persistor == nil {
			persistor = &LocalMessagePersister{
				queue:  make(chan Envelope, 100),
				logger: nil,
			}
		}
	})

	return persistor
}

func SetGlobalPersistor(p *LocalMessagePersister) {
	persistor = p
}

type MessagePersister interface {
	Persist(message Envelope) error
}

type LocalMessagePersister struct {
	logger *zap.Logger
	queue  chan Envelope
}

func NewLocalMessagePersister(l *zap.Logger) *LocalMessagePersister {
	return &LocalMessagePersister{
		queue: make(chan Envelope, 100),
	}
}

func (lp *LocalMessagePersister) Persist(message Envelope) error {
	lp.queue <- message

	return nil
}

func (lp *LocalMessagePersister) Channel() chan Envelope {
	return lp.queue
}
