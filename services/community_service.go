package services

import (
	"github.com/vicebe/following-service/data"
	"log"
)

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

func (cs *CommunityService) FollowCommunity(cID string, uID string) error {

	community, err := cs.cr.FindBy("external_id", cID)

	if err != nil {
		return err
	}

	user, err := cs.ur.FindBy("external_id", uID)

	if err != nil {
		return nil
	}

	err = cs.cr.FollowCommunity(community, user)

	if err != nil {
		return err
	}

	return nil
}

func (cs *CommunityService) UnfollowCommunity(cID string, uID string) error {

	community, err := cs.cr.FindBy("external_id", cID)

	if err != nil {
		return err
	}

	user, err := cs.ur.FindBy("external_id", uID)

	if err != nil {
		return nil
	}

	err = cs.cr.UnfollowCommunity(community, user)

	if err != nil {
		return err
	}

	return nil
}
