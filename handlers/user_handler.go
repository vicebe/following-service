package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/services"
	"log"
	"net/http"
)

type UserHandler struct {
	l           *log.Logger
	userService *services.UserService
}

func NewUserHandler(
	l *log.Logger,
	userService *services.UserService,
) *UserHandler {
	return &UserHandler{
		l:           l,
		userService: userService,
	}
}

// GetFollowers is a GET handler that returns all the followers of a user
func (uh *UserHandler) GetFollowers(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")
	uId := chi.URLParam(r, "userId")

	followers, err := uh.userService.GetUserFollowers(uId)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJson(&SimpleResponse{Message: err.Error()}, rw)
		return
	}

	data.ToJson(&FollowersResponse{Followers: followers}, rw)
}

// FollowUser is POST handler that handles request for a user to follow another
// user
func (uh *UserHandler) FollowUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	uID := chi.URLParam(r, "userId")
	fID := chi.URLParam(r, "toFollowId")

	err := uh.userService.FollowUser(uID, fID)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJson(&SimpleResponse{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// UnfollowUser is a DELETE handler that handles requests to unfollow users
func (uh *UserHandler) UnfollowUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	uID := chi.URLParam(r, "userId")
	fID := chi.URLParam(r, "toUnfollowId")

	err := uh.userService.UnfollowUser(uID, fID)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJson(&SimpleResponse{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
