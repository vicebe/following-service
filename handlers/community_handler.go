package handlers

import (
	"github.com/go-chi/chi/v5"
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
func (ch *CommunityHandler) FollowCommunity(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	uID := chi.URLParam(r, "userId")
	cID := chi.URLParam(r, "communityId")

	err := ch.communityService.FollowCommunity(cID, uID)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJson(&SimpleResponse{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
