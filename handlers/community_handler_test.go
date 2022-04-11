package handlers_test

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/vicebe/following-service/handlers"
	"github.com/vicebe/following-service/handlers/test_utils"
	"github.com/vicebe/following-service/services"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type middleware func(http.Handler) http.Handler

func TestCommunityHandler_FollowCommunity(t *testing.T) {
	// responses

	internalError := handlers.MakeInternalErrorResponse()
	internalErrorResponse, _ := json.Marshal(&internalError)

	// tests cases
	tests := []struct {
		name             string
		l                *log.Logger
		communityService services.CommunityServiceI
		method           string
		statusCode       int
		responseBody     string
		middlewares      []middleware
	}{
		{
			name: "test follow community success",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			communityService: test_utils.CommunityServiceFollowCommunityMock{},
			method:           http.MethodPost,
			statusCode:       http.StatusNoContent,
			responseBody:     "",
			middlewares: []middleware{
				test_utils.AddCommunityToRequestContext(&test_utils.CommunityOne),
				test_utils.AddUserToRequestContext(&test_utils.UserOne),
			},
		},
		{
			name: "test community not passed in context",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			communityService: test_utils.CommunityServiceFollowCommunityMock{},
			method:           http.MethodPost,
			statusCode:       http.StatusInternalServerError,
			responseBody:     string(internalErrorResponse),
			middlewares: []middleware{
				test_utils.IdentityMiddleware,
				test_utils.AddUserToRequestContext(&test_utils.UserOne),
			},
		},
		{
			name: "test user not passed in context",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			communityService: test_utils.CommunityServiceFollowCommunityMock{},
			method:           http.MethodPost,
			statusCode:       http.StatusInternalServerError,
			responseBody:     string(internalErrorResponse),
			middlewares: []middleware{
				test_utils.AddCommunityToRequestContext(&test_utils.CommunityOne),
				test_utils.IdentityMiddleware,
			},
		},
		{
			name: "test community to follow error",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			communityService: test_utils.CommunityServiceFollowCommunityErrorMock{},
			method:           http.MethodPost,
			statusCode:       http.StatusInternalServerError,
			responseBody:     string(internalErrorResponse),
			middlewares: []middleware{
				test_utils.AddCommunityToRequestContext(&test_utils.CommunityOne),
				test_utils.AddUserToRequestContext(&test_utils.UserOne),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				tt.method,
				"/follow-user",
				nil,
			)

			rr := httptest.NewRecorder()

			ch := handlers.NewCommunityHandler(tt.l, tt.communityService)

			r := chi.NewRouter()

			for _, mid := range tt.middlewares {
				r.Use(mid)
			}

			r.Post("/follow-user", ch.FollowCommunity)

			r.ServeHTTP(rr, request)

			if rr.Code != tt.statusCode {
				t.Errorf(
					"expected http status code %d got %d",
					tt.statusCode,
					rr.Code,
				)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.responseBody {
				t.Errorf(
					"expected response '%s' got '%s'",
					tt.responseBody,
					rr.Body.String(),
				)
			}
		})
	}

}

func TestCommunityHandler_UnfollowCommunity(t *testing.T) {
	// responses
	internalError := handlers.MakeInternalErrorResponse()
	internalErrorResponse, _ := json.Marshal(&internalError)

	// tests cases
	tests := []struct {
		name             string
		l                *log.Logger
		communityService services.CommunityServiceI
		method           string
		statusCode       int
		responseBody     string
		middlewares      []middleware
	}{
		{
			name: "test unfollow community success",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			communityService: test_utils.CommunityServiceUnfollowCommunityMock{},
			method:           http.MethodDelete,
			statusCode:       http.StatusNoContent,
			responseBody:     "",
			middlewares: []middleware{
				test_utils.AddCommunityToRequestContext(&test_utils.CommunityOne),
				test_utils.AddUserToRequestContext(&test_utils.UserOne),
			},
		},
		{
			name: "test community not passed in context",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			communityService: test_utils.CommunityServiceUnfollowCommunityMock{},
			method:           http.MethodDelete,
			statusCode:       http.StatusInternalServerError,
			responseBody:     string(internalErrorResponse),
			middlewares: []middleware{
				test_utils.IdentityMiddleware,
				test_utils.AddUserToRequestContext(&test_utils.UserOne),
			},
		},
		{
			name: "test user not passed in context",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			communityService: test_utils.CommunityServiceUnfollowCommunityMock{},
			method:           http.MethodDelete,
			statusCode:       http.StatusInternalServerError,
			responseBody:     string(internalErrorResponse),
			middlewares: []middleware{
				test_utils.AddCommunityToRequestContext(&test_utils.CommunityOne),
				test_utils.IdentityMiddleware,
			},
		},
		{
			name: "test community to unfollow error",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			communityService: test_utils.CommunityServiceUnfollowCommunityErrorMock{},
			method:           http.MethodDelete,
			statusCode:       http.StatusInternalServerError,
			responseBody:     string(internalErrorResponse),
			middlewares: []middleware{
				test_utils.AddCommunityToRequestContext(&test_utils.CommunityOne),
				test_utils.AddUserToRequestContext(&test_utils.UserOne),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				tt.method,
				"/unfollow-user",
				nil,
			)

			rr := httptest.NewRecorder()

			ch := handlers.NewCommunityHandler(tt.l, tt.communityService)

			r := chi.NewRouter()

			for _, mid := range tt.middlewares {
				r.Use(mid)
			}

			r.Delete("/unfollow-user", ch.UnfollowCommunity)

			r.ServeHTTP(rr, request)

			if rr.Code != tt.statusCode {
				t.Errorf(
					"expected http status code %d got %d",
					tt.statusCode,
					rr.Code,
				)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.responseBody {
				t.Errorf(
					"expected response '%s' got '%s'",
					tt.responseBody,
					rr.Body.String(),
				)
			}
		})
	}
}

func TestCommunityHandler_GetCommunityFollowers(t *testing.T) {
	// responses
	followersResponse, _ := json.Marshal(
		&handlers.FollowersResponse{Followers: test_utils.FollowersList},
	)

	internalError := handlers.MakeInternalErrorResponse()
	internalErrorResponse, _ := json.Marshal(&internalError)

	// tests cases
	tests := []struct {
		name             string
		l                *log.Logger
		communityService services.CommunityServiceI
		method           string
		statusCode       int
		responseBody     string
		middleware       func(http.Handler) http.Handler
	}{
		{
			name: "test get followers success",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			communityService: test_utils.CommunityServiceGetFollowersMock{},
			method:           http.MethodGet,
			statusCode:       http.StatusOK,
			responseBody:     string(followersResponse),
			middleware:       test_utils.AddCommunityToRequestContext(&test_utils.CommunityOne),
		},
		{
			name: "test community not passed in context",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			communityService: test_utils.CommunityServiceGetFollowersMock{},
			method:           http.MethodGet,
			statusCode:       http.StatusInternalServerError,
			responseBody:     string(internalErrorResponse),
			middleware:       test_utils.IdentityMiddleware,
		},
		{
			name: "test get followers errors",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			communityService: test_utils.CommunityServiceGetFollowersErrorMock{},
			method:           http.MethodGet,
			statusCode:       http.StatusInternalServerError,
			responseBody:     string(internalErrorResponse),
			middleware:       test_utils.AddCommunityToRequestContext(&test_utils.CommunityOne),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				tt.method,
				"/get-followers",
				nil,
			)

			rr := httptest.NewRecorder()

			ch := handlers.NewCommunityHandler(tt.l, tt.communityService)

			r := chi.NewRouter()

			r.Use(tt.middleware)

			r.Get("/get-followers", ch.GetCommunityFollowers)

			r.ServeHTTP(rr, request)

			if rr.Code != tt.statusCode {
				t.Errorf(
					"expected http status code %d got %d",
					tt.statusCode,
					rr.Code,
				)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.responseBody {
				t.Errorf(
					"expected response '%s' got '%s'",
					tt.responseBody,
					rr.Body.String(),
				)
			}
		})
	}
}
