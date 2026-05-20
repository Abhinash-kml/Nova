package realtime

type RealtimeBroker interface {
	Initialize() error
	Publish(string, Envelope) error
	Subscribe(string)
	Unsubscribe(string)
	ListenToSubscriptions() <-chan Envelope
}
