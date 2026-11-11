package messagebroker

import "sync"

var (
	globalInstance MessageBroker
	once           sync.Once
)

type MessageBroker interface {
	Initialize()
	Produce(channel string, payload map[string]any)
	Consume(channel string) <-chan map[string]any
}

func SetGlobal(mb MessageBroker) {
	once.Do(func() {
		globalInstance = mb
	})
}

func GetGlobal() MessageBroker {
	if globalInstance == nil {
		panic("Global events manager is nil")
	}
	return globalInstance
}
