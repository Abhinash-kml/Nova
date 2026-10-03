package realtime

type MessagePersister interface {
	Persist(message Envelope) error
}
