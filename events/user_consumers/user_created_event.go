package userconsumers

import (
	"encoding/json"
	"log"

	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/services"
)

type UserMessage struct {
	ExternalID string `json:"external_id"`
}
type UserCreatedEvent struct {
	User UserMessage `json:"user"`
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

	user := &data.User{
		ExternalID: userCreatedEvent.User.ExternalID,
	}

	if err := ucc.userService.CreateUser(user); err != nil {
		return err
	}

	return nil
}
