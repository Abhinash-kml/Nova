package events

import "sync"

var (
	globalInstance EventsManager
	once           sync.Once
)

type EventsManager interface {
	Initialize()
	Produce(channel string, payload map[string]any)
	Consume(channel string) <-chan map[string]any
}

func SetGlobalEventsManager(em EventsManager) {
	once.Do(func() {
		globalInstance = em
	})
}

func GetGlobalEventsManager() EventsManager {
	if globalInstance == nil {
		panic("Global events manager is nil")
	}
	return globalInstance
}
