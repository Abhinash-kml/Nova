package realtime

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

var (
	persistorOnce   sync.Once
	persistorGlobal *LocalMessagePersister
)

func GetLocalMessagePersistor() *LocalMessagePersister {
	if persistorGlobal == nil {
		panic("global instance of local message persistor is nil")
	}

	return persistorGlobal
}

func SetGlobalPersistor(p *LocalMessagePersister) error {
	if p == nil {
		return errors.New("cannot set local message persistor with nil pointer")
	}
	persistorOnce.Do(func() {
		persistorGlobal = p
	})
	return nil
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
