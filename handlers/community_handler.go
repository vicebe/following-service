package handlers

import (
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/services"
	"log"
	"net/http"
)

type CommunityHandler struct {
	l                *log.Logger
	communityService *services.CommunityService
}

func NewCommunityHandler(
	l *log.Logger,
	communityService *services.CommunityService,
) *CommunityHandler {
	return &CommunityHandler{
		l:                l,
		communityService: communityService,
	}
}

// FollowCommunity is a POST handler that handles requests for a user to follow
// a community
func (ch *CommunityHandler) FollowCommunity(
	rw http.ResponseWriter,
	r *http.Request,
) {
	rw.Header().Add("Content-Type", "application/json")

	community, ok := r.Context().Value("community").(*data.Community)

	if !ok {
		ch.l.Print("[ERROR] user not passed in context")
		SetInternalErrorResponse(rw, ch.l)
		return
	}

	user, ok := r.Context().Value("user").(*data.User)

	if !ok {
		ch.l.Printf("[ERROR] user not passed in context")
		SetInternalErrorResponse(rw, ch.l)
		return
	}

	err := ch.communityService.FollowCommunity(community, user)

	if err != nil {
		SetInternalErrorResponse(rw, ch.l)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// UnfollowCommunity is a DELETE handler that handles requests for a user to
// unfollow a community
func (ch *CommunityHandler) UnfollowCommunity(
	rw http.ResponseWriter,
	r *http.Request,
) {

	rw.Header().Add("Content-Type", "application/json")

	community, ok := r.Context().Value("community").(*data.Community)

	if !ok {
		ch.l.Print("[ERROR] user not passed in context")
		SetInternalErrorResponse(rw, ch.l)
		return
	}

	user, ok := r.Context().Value("user").(*data.User)

	if !ok {
		ch.l.Print("[ERROR] user not passed in context")
		SetInternalErrorResponse(rw, ch.l)
		return
	}

	err := ch.communityService.UnfollowCommunity(community, user)

	if err != nil {
		SetInternalErrorResponse(rw, ch.l)
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}

func (ch *CommunityHandler) GetCommunityFollowers(
	rw http.ResponseWriter, r *http.Request,
) {
	rw.Header().Add("Content-Type", "application/json")

	community, ok := r.Context().Value("community").(*data.Community)

	if !ok {
		ch.l.Print("[ERROR] user not passed in context")
		SetInternalErrorResponse(rw, ch.l)
		return
	}

	followers, err := ch.communityService.GetCommunityFollowers(community)

	if err != nil {
		SetInternalErrorResponse(rw, ch.l)
	}

	if err = data.ToJson(
		&FollowersResponse{Followers: followers},
		rw,
	); err != nil {
		ch.l.Print("[ERROR] ", err)
		SetInternalErrorResponse(rw, ch.l)
	}
}
