package events

import (
	"context"
	"errors"
	"log"
)

// KafkaConsumerClient is client for that connects and talks to a kafka broker.
type KafkaConsumerClient interface {

	// ReadMessage reads a message from the kafka broker.
	ReadMessage(context.Context) ([]byte, error)

	// Close closes the kafka broker.
	Close() error
}

// KafkaConsumer implements the consumer interface to listen and consume
// messages from a kafka broker.
type KafkaConsumer struct {
	reader       KafkaConsumerClient
	Topic        string
	logger       *log.Logger
	consumeEvent ConsumerFunc
	stopCtx      context.Context
	stopFunc     context.CancelFunc
}

func NewKafkaConsumer(
	topic string,
	reader KafkaConsumerClient,
	logger *log.Logger,
	consumerFunc ConsumerFunc,
) *KafkaConsumer {
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &KafkaConsumer{
		reader:       reader,
		Topic:        topic,
		logger:       logger,
		consumeEvent: consumerFunc,
		stopCtx:      ctx,
		stopFunc:     cancelFunc,
	}
}

func (c *KafkaConsumer) StartConsumer() error {

	c.logger.Printf(
		"[INFO]: Starting consumer for topic %s\n",
		c.Topic,
	)

	go c.RunConsumer()

	return nil
}

func (c *KafkaConsumer) RunConsumer() {
	defer func() {
		if err := c.reader.Close(); err != nil {
			c.logger.Print("[ERROR]: ", err)
		}
	}()

	for {

		m, err := c.reader.ReadMessage(c.stopCtx)

		if err != nil {

			if errors.Is(err, context.Canceled) {
				break
			}

			c.logger.Print("[ERROR]: ", err)

		} else {
			_ = c.ConsumeEvent(m)
		}

	}
}

func (c *KafkaConsumer) ConsumeEvent(message []byte) error {
	err := c.consumeEvent(message)
	if err != nil {
		c.logger.Print("[ERROR]: ", err)
		return err
	}

	return nil
}

func (c *KafkaConsumer) StopConsumer() error {
	c.logger.Print(
		"[INFO]: Stoping consumer of topic ",
		c.Topic,
	)
	c.stopFunc()
	return nil
}
