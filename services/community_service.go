package services

import (
	"log"

	"github.com/vicebe/following-service/data"
)

type CommunityServiceI interface {
	GetCommunity(cID string) (*data.Community, error)
	FollowCommunity(community *data.Community, user *data.User) error
	UnfollowCommunity(community *data.Community, user *data.User) error
	GetCommunityFollowers(community *data.Community) ([]data.User, error)
}

// CommunityService is a service that handles common business logic for
// communities.
type CommunityService struct {
	l  *log.Logger
	cr data.CommunityRepository
	ur data.UserRepository
}

func NewCommunityService(
	l *log.Logger,
	cr data.CommunityRepository,
	ur data.UserRepository,
) *CommunityService {
	return &CommunityService{l, cr, ur}
}

// GetCommunity retrieves community given its ID.
func (cs *CommunityService) GetCommunity(cID string) (*data.Community, error) {

	community, err := cs.cr.FindBy("external_id", cID)

	if err != nil {
		return nil, err
	}

	return community, nil
}

func (cs *CommunityService) FollowCommunity(
	community *data.Community,
	user *data.User,
) error {

	err := cs.cr.FollowCommunity(community, user)

	if err != nil {
		return err
	}

	return nil
}

func (cs *CommunityService) UnfollowCommunity(
	community *data.Community,
	user *data.User,
) error {

	err := cs.cr.UnfollowCommunity(community, user)

	if err != nil {
		return err
	}

	return nil
}

func (cs *CommunityService) GetCommunityFollowers(
	community *data.Community,
) ([]data.User, error) {

	followers, err := cs.cr.GetCommunityFollowers(community)

	if err != nil {
		return nil, err
	}

	return followers, nil

}

func (cs *CommunityService) CreateCommunity(community *data.Community) error {
	_, err := cs.cr.FindBy("external_id", community.ExternalID)

	if err == data.ErrorCommunityNotFound {
		if err := cs.cr.Create(community); err != nil {
			return err
		}

		return nil
	}

	return err
}
