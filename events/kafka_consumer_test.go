package events_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/vicebe/following-service/events"
)

type KafkaConsumerClientMock struct {
	ReadMessageCalled bool
}

func NewKafkaConsumerClientMock() *KafkaConsumerClientMock {
	return &KafkaConsumerClientMock{ReadMessageCalled: false}
}

func (k *KafkaConsumerClientMock) ReadMessage(_ context.Context) ([]byte, error) {
	if k.ReadMessageCalled {
		return nil, context.Canceled
	}

	k.ReadMessageCalled = true

	return []byte("message"), nil
}

func (k *KafkaConsumerClientMock) Close() error {
	return nil
}

func TestKafkaConsumer_StartConsumer(t *testing.T) {

	t.Run("test start consumer", func(t *testing.T) {
		kafkaClient := NewKafkaConsumerClientMock()
		reader := events.NewKafkaConsumer(
			"test-topic",
			kafkaClient,
			log.New(os.Stdout, "test", log.LstdFlags),
			func(message []byte) error {
				return nil
			},
		)

		reader.RunConsumer()

		if !kafkaClient.ReadMessageCalled {
			t.Error("ReadMessage Was not called")
		}

	})

}
