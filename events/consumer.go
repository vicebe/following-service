package events

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/services"
	"log"
)

type Consumer interface {
	StartConsumer() error
	ConsumeEvent(message []byte) error
	StopConsumer() error
}

type ConsumerFunc func([]byte) error

type KafkaConsumer struct {
	reader       *kafka.Reader
	logger       *log.Logger
	stopChan     chan bool
	consumeEvent ConsumerFunc
}

func NewKafkaConsumer(
	config kafka.ReaderConfig,
	logger *log.Logger,
	consumerFunc ConsumerFunc,
) *KafkaConsumer {
	return &KafkaConsumer{
		reader:       kafka.NewReader(config),
		logger:       logger,
		consumeEvent: consumerFunc,
		stopChan:     nil,
	}
}

func (c KafkaConsumer) StartConsumer() error {

	c.stopChan = make(chan bool)

	go func() {

		defer func() {
			if err := c.reader.Close(); err != nil {
				c.logger.Print("[ERROR] ", err)
			}
		}()

		for {

			select {

			case <-c.stopChan:
				break

			default:
				m, err := c.reader.ReadMessage(context.Background())

				if err != nil {
					c.logger.Print("[ERROR] ", err)
					continue
				}

				_ = c.ConsumeEvent(m.Value)
			}
		}

	}()

	return nil
}

func (c KafkaConsumer) ConsumeEvent(message []byte) error {
	err := c.consumeEvent(message)
	if err != nil {
		c.logger.Print("[ERROR] ", err)
		return err
	}

	return nil
}

func (c KafkaConsumer) StopConsumer() error {
	c.stopChan <- true
	return nil
}

// Events

type UserCreatedEvent struct {
	User data.User
}

type UserCreatedConsumer struct {
	logger      *log.Logger
	userService *services.UserService
}

func NewUserCreatedConsumer(
	logger *log.Logger,
	userService *services.UserService,
) *UserCreatedConsumer {
	return &UserCreatedConsumer{
		logger:      logger,
		userService: userService,
	}
}

func (ucc *UserCreatedConsumer) UserCreatedEventHandler(message []byte) error {
	var userCreatedEvent UserCreatedEvent

	if err := json.Unmarshal(message, &userCreatedEvent); err != nil {
		ucc.logger.Print("[ERROR] ", err)
		return err
	}

	if err := ucc.userService.CreateUser(&userCreatedEvent.User); err != nil {
		return err
	}

	return nil
}

type CommunityCreatedEvent struct {
	Community data.Community
}

type CommunityCreatedConsumer struct {
	logger           *log.Logger
	communityService *services.CommunityService
}

func NewCommunityCreatedConsumer(
	logger *log.Logger,
	communityService *services.CommunityService,
) *CommunityCreatedConsumer {
	return &CommunityCreatedConsumer{
		logger:           logger,
		communityService: communityService,
	}
}

func (ucc *CommunityCreatedConsumer) CommunityCreatedEventHandler(
	message []byte,
) error {
	var communityCreatedEvent CommunityCreatedEvent

	if err := json.Unmarshal(message, &communityCreatedEvent); err != nil {
		ucc.logger.Print("[ERROR] ", err)
		return err
	}

	if err := ucc.communityService.CreateCommunity(
		&communityCreatedEvent.Community,
	); err != nil {
		return err
	}

	return nil
}
