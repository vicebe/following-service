package events

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer implements the producer interface to produce messages to a
// kafka broker.
type KafkaProducer struct {
	writer *kafka.Writer
	logger *log.Logger
}

func NewKafkaProducer(
	config kafka.WriterConfig,
	logger *log.Logger,
) *KafkaProducer {
	return &KafkaProducer{
		writer: kafka.NewWriter(config),
		logger: logger,
	}
}
func (kp *KafkaProducer) StartProducer() error {
	panic("not implemented") // TODO: Implement
}

// ProduceEvent sends an event to the broker.
func (kp *KafkaProducer) ProduceEvent(message []byte) error {
	if err := kp.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: message,
		},
	); err != nil {
		kp.logger.Print("[ERROR]: error sending message to broker")
		kp.logger.Print("[ERROR]: ", err)
		return err
	}

	return nil
}

func (kp *KafkaProducer) StopProducer() error {
	if err := kp.writer.Close(); err != nil {
		kp.logger.Print("[ERROR]: Error stopping producer")
		kp.logger.Print("[ERROR]: ", err)
		return err
	}

	return nil
}
