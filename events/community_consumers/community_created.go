package communityconsumers

import (
	"encoding/json"
	"log"

	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/services"
)

type CommunityMessage struct {
	ExternalID string `json:"external_id"`
}

type CommunityCreatedEvent struct {
	Community CommunityMessage `json:"community"`
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

	community := &data.Community{
		ExternalID: communityCreatedEvent.Community.ExternalID,
	}

	if err := ucc.communityService.CreateCommunity(community); err != nil {
		return err
	}

	return nil
}
