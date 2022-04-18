package events

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaGoConsumerClient implements the KafkaConsumerClient interface using the
// kafka-go library.
type KafkaGoConsumerClient struct {
	reader *kafka.Reader
}

func NewKafkaGoConsumer(reader *kafka.Reader) *KafkaGoConsumerClient {
	return &KafkaGoConsumerClient{
		reader: reader,
	}
}

func (k *KafkaGoConsumerClient) ReadMessage(ctx context.Context) ([]byte, error) {
	m, err := k.reader.ReadMessage(ctx)

	if err != nil {
		return nil, err
	}

	return m.Value, nil
}

func (k *KafkaGoConsumerClient) Close() error {
	return k.Close()
}
