package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/services"
)

// Handler is a handler for the service routes.
type Handler struct {
	l  *log.Logger
	as *services.AppService
}

func NewHandler(l *log.Logger, as *services.AppService) *Handler {
	return &Handler{l, as}
}

// GetFollowers is a GET handler that returns all the followers of a user
func (sh *Handler) GetFollowers(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")

	uId := chi.URLParam(r, "userId")

	sh.l.Printf("[DEBUG] finding user %v\n", uId)

	followers, err := sh.as.GetFollowers(uId)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJson(&SimpleResponse{Message: err.Error()}, rw)
		return
	}

	data.ToJson(&FollowersResponse{Followers: followers}, rw)
}

// FollowUser is POST handler that handles request for a user to follow another
// user
func (sh *Handler) FollowUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	uId := chi.URLParam(r, "userId")
	fId := chi.URLParam(r, "toFollowId")

	sh.l.Printf("[DEBUG] user %v wants to follow %v\n", uId, fId)

	err := sh.as.FollowUser(uId, fId)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJson(&SimpleResponse{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// UnfollowUser is a DELETE handler that handles resquests to unfollow users
func (sh *Handler) UnfollowUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	uId := chi.URLParam(r, "userId")
	fId := chi.URLParam(r, "toUnfollowId")

	err := sh.as.UnfollowUser(uId, fId)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJson(&SimpleResponse{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// FollowCommunity is a POST handler that handles requests for a user to follow a
// community
func (sh *Handler) FollowCommunity(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	uId := chi.URLParam(r, "userId")
	cId := chi.URLParam(r, "communityId")

	err := sh.as.FollowCommunity(uId, cId)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJson(&SimpleResponse{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
