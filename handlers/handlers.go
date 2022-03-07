package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/services"
)

type SimpleResponse struct {
	Message string `json:"message"`
}

type FollowersResponse struct {
	Followers []string `json:"followers"`
}

type ServiceHandler struct {
	l  *log.Logger
	us *services.UserService
}

func NewServiceHandler(
	l *log.Logger, us *services.UserService,
) *ServiceHandler {
	return &ServiceHandler{l, us}
}

func (sh *ServiceHandler) GetFollowers(
	rw http.ResponseWriter, r *http.Request,
) {

	rw.Header().Add("Content-Type", "application/json")

	uId := chi.URLParam(r, "userId")

	sh.l.Printf("[DEBUG] finding user %v\n", uId)

	followers, err := sh.us.GetFollowers(uId)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJson(&SimpleResponse{Message: err.Error()}, rw)
		return
	}

	data.ToJson(&FollowersResponse{Followers: followers}, rw)
}

func (sh *ServiceHandler) FollowUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	uId := chi.URLParam(r, "userId")
	fId := chi.URLParam(r, "toFollowId")

	sh.l.Printf("[DEBUG] user %v wants to follow %v\n", uId, fId)

	err := sh.us.FollowUser(uId, fId)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJson(&SimpleResponse{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
