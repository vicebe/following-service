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
	uID := chi.URLParam(r, "userID")

	followers, err := uh.userService.GetUserFollowers(uID)

	if err != nil {
		var httpStatus int
		var response SimpleResponse

		if err == data.ErrorUserNotFound {
			httpStatus = http.StatusNotFound
			response = SimpleResponse{Message: "User not found"}
		} else {
			httpStatus = http.StatusInternalServerError
			response = SimpleResponse{Message: "Something went wrong"}

			uh.l.Print("[Error] ", err)
		}

		rw.WriteHeader(httpStatus)

		if err := data.ToJson(&response, rw); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			uh.l.Print("[Error] ", err)
		}

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

	uID := chi.URLParam(r, "userID")
	fID := chi.URLParam(r, "followerID")

	err := uh.userService.FollowUser(fID, uID)

	if err != nil {
		var httpStatus int
		var response SimpleResponse

		if err == data.ErrorUserNotFound {
			httpStatus = http.StatusNotFound
			response = SimpleResponse{Message: "User not found"}
		} else {
			httpStatus = http.StatusInternalServerError
			response = SimpleResponse{Message: "Something went wrong"}

			uh.l.Print("[Error] ", err)
		}

		rw.WriteHeader(httpStatus)

		if err := data.ToJson(&response, rw); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			uh.l.Print("[Error] ", err)
		}

		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// UnfollowUser is a DELETE handler that handles requests to unfollow users
func (uh *UserHandler) UnfollowUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	uID := chi.URLParam(r, "userID")
	fID := chi.URLParam(r, "followerID")

	err := uh.userService.UnfollowUser(fID, uID)

	if err != nil {
		var httpStatus int
		var response SimpleResponse

		if err == data.ErrorUserNotFound {
			httpStatus = http.StatusNotFound
			response = SimpleResponse{Message: "User not found"}
		} else {
			httpStatus = http.StatusInternalServerError
			response = SimpleResponse{Message: "Something went wrong"}

			uh.l.Print("[Error] ", err)
		}

		rw.WriteHeader(httpStatus)

		if err := data.ToJson(&response, rw); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			uh.l.Print("[Error] ", err)
		}

		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// GetCommunities is a GET handler that returns all the communities that the
// user follows
func (uh *UserHandler) GetCommunities(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")
	uID := chi.URLParam(r, "userID")

	communities, err := uh.userService.GetUserCommunities(uID)

	if err != nil {
		var httpStatus int
		var response SimpleResponse

		if err == data.ErrorUserNotFound {
			httpStatus = http.StatusNotFound
			response = SimpleResponse{Message: "User not found"}
		} else {
			httpStatus = http.StatusInternalServerError
			response = SimpleResponse{Message: "Something went wrong"}

			uh.l.Print("[Error] ", err)
		}

		rw.WriteHeader(httpStatus)

		if err := data.ToJson(&response, rw); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			uh.l.Print("[Error] ", err)
		}

		return
	}

	response := &CommunitiesResponse{Communities: communities}
	if err := data.ToJson(&response, rw); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		uh.l.Print("[Error] ", err)
	}
}
