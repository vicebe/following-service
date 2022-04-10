package handlers

import (
	"github.com/vicebe/following-service/data"
	"log"
	"net/http"
)

type SimpleResponse struct {
	Message string `json:"message"`
}

func MakeInternalErrorResponse() SimpleResponse {
	return SimpleResponse{Message: "Something went wrong"}
}

// SetInternalErrorResponse sets code status to 500 and error response.
func SetInternalErrorResponse(rw http.ResponseWriter, logger *log.Logger) {
	rw.WriteHeader(http.StatusInternalServerError)

	response := MakeInternalErrorResponse()
	if err := data.ToJson(&response, rw); err != nil {
		logger.Print("[ERROR] ", err)
	}

}

type FollowersResponse struct {
	Followers []data.User `json:"followers"`
}

type CommunitiesResponse struct {
	Communities []data.Community `json:"communities"`
}
