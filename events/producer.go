package events

// Producer is an interface that produces messages to a broker.
type Producer interface {
	StartProducer() error

	// ProduceEvent sends an event to the broker.
	ProduceEvent(message []byte) error

	StopProducer() error
}
