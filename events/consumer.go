package events

// Consumer is an interface that consumes messages from a broker.
type Consumer interface {

	// StartConsumer starts a consumer to listen to borker for messages.
	StartConsumer() error

	// ConsumeEvent handles an event from the broker.
	ConsumeEvent(message []byte) error

	// StopConsumer gracefully stops consumer.
	StopConsumer() error
}

// ConsumerFunc is a function that recieves a message from the broker.
type ConsumerFunc func([]byte) error
