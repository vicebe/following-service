package events

import (
	"context"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaConsumer implements the consumer interface to listen and consume
// messages from a kafka broker.
type KafkaConsumer struct {
	reader       *kafka.Reader
	logger       *log.Logger
	consumeEvent ConsumerFunc
	stopCtx      context.Context
	stopFunc     context.CancelFunc
}

func NewKafkaConsumer(
	config kafka.ReaderConfig,
	logger *log.Logger,
	consumerFunc ConsumerFunc,
) *KafkaConsumer {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &KafkaConsumer{
		reader:       kafka.NewReader(config),
		logger:       logger,
		consumeEvent: consumerFunc,
		stopCtx:      ctx,
		stopFunc:     cancelFunc,
	}
}

func (c *KafkaConsumer) StartConsumer() error {

	c.logger.Printf(
		"[INFO]: Starting consumer for topic %s\n",
		c.reader.Config().Topic,
	)

	go func() {

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
				_ = c.ConsumeEvent(m.Value)
			}

		}

	}()

	return nil
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
		c.reader.Config().Topic,
	)
	c.stopFunc()
	return nil
}
