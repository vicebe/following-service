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
	userService services.UserServiceI
}

func NewUserHandler(
	l *log.Logger,
	userService services.UserServiceI,
) *UserHandler {
	return &UserHandler{
		l:           l,
		userService: userService,
	}
}

// GetFollowers is a GET handler that returns all the followers of a user
func (uh *UserHandler) GetFollowers(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")

	user, ok := r.Context().Value("user").(*data.User)

	if !ok {
		uh.l.Printf("[ERROR] user not passed in context")
		SetInternalErrorResponse(rw, uh.l)
		return
	}

	followers, err := uh.userService.GetUserFollowers(user)

	if err != nil {
		SetInternalErrorResponse(rw, uh.l)
		return
	}

	response := &FollowersResponse{Followers: followers}
	if err := data.ToJson(&response, rw); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		uh.l.Print("[Error] ", err)
	}
}

// FollowUser is POST handler that handles request for a user to follow another
// user
func (uh *UserHandler) FollowUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	user, ok := r.Context().Value("user").(*data.User)

	if !ok {
		uh.l.Printf("[ERROR] user not passed in context")
		SetInternalErrorResponse(rw, uh.l)
		return
	}

	followerID := chi.URLParam(r, "followerID")

	follower, err := uh.userService.GetUser(followerID)

	if err != nil {
		var errorStatus int
		var errorResponse SimpleResponse

		switch err {
		case data.ErrorUserNotFound:
			errorStatus = http.StatusNotFound
			errorResponse = SimpleResponse{Message: "User not found"}
		default:
			errorStatus = http.StatusInternalServerError
			errorResponse = MakeInternalErrorResponse()
		}

		rw.WriteHeader(errorStatus)
		if err := data.ToJson(&errorResponse, rw); err != nil {
			SetInternalErrorResponse(rw, uh.l)
		}

		return
	}

	err = uh.userService.FollowUser(follower, user)

	if err != nil {
		SetInternalErrorResponse(rw, uh.l)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// UnfollowUser is a DELETE handler that handles requests to unfollow users
func (uh *UserHandler) UnfollowUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	user, ok := r.Context().Value("user").(*data.User)

	if !ok {
		uh.l.Printf("[ERROR] user not passed in context")
		SetInternalErrorResponse(rw, uh.l)
		return
	}

	followerID := chi.URLParam(r, "followerID")

	follower, err := uh.userService.GetUser(followerID)

	if err != nil {
		var errorStatus int
		var errorResponse SimpleResponse

		switch err {
		case data.ErrorUserNotFound:
			errorStatus = http.StatusNotFound
			errorResponse = SimpleResponse{Message: "User not found"}
		default:
			errorStatus = http.StatusInternalServerError
			errorResponse = MakeInternalErrorResponse()
		}

		rw.WriteHeader(errorStatus)
		if err := data.ToJson(&errorResponse, rw); err != nil {
			SetInternalErrorResponse(rw, uh.l)
		}

		return
	}

	err = uh.userService.UnfollowUser(follower, user)

	if err != nil {
		SetInternalErrorResponse(rw, uh.l)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// GetCommunities is a GET handler that returns all the communities that the
// user follows
func (uh *UserHandler) GetCommunities(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")
	user, ok := r.Context().Value("user").(*data.User)

	if !ok {
		uh.l.Printf("[ERROR] user not passed in context")
		SetInternalErrorResponse(rw, uh.l)
		return
	}

	communities, err := uh.userService.GetUserCommunities(user)

	if err != nil {
		SetInternalErrorResponse(rw, uh.l)
		return
	}

	response := &CommunitiesResponse{Communities: communities}
	if err := data.ToJson(&response, rw); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		uh.l.Print("[Error] ", err)
	}
}
