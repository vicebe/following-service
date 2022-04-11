package handlers_test

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/handlers/test_utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/vicebe/following-service/handlers"
	"github.com/vicebe/following-service/services"
)

func TestUserHandler_FollowUser(t *testing.T) {
	// responses
	userNotFoundResponse, _ := json.Marshal(
		&handlers.SimpleResponse{Message: "User not found"},
	)

	internalError := handlers.MakeInternalErrorResponse()
	internalErrorResponse, _ := json.Marshal(&internalError)

	// tests cases
	tests := []struct {
		name         string
		l            *log.Logger
		userService  services.UserServiceI
		method       string
		statusCode   int
		responseBody string
		middleware   func(http.Handler) http.Handler
	}{
		{
			name: "test follow user success",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceFollowUserMock{},
			method:       http.MethodPost,
			statusCode:   http.StatusNoContent,
			responseBody: "",
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
		},
		{
			name: "test user not passed in context",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceFollowUserMock{},
			method:       http.MethodPost,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.IdentityMiddleware,
		},
		{
			name: "test user to follow not found",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceFollowUserNotFoundMock{},
			method:       http.MethodPost,
			statusCode:   http.StatusNotFound,
			responseBody: string(userNotFoundResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
		},
		{
			name: "test get user to follow error",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceFollowUserGetUserErrorMock{},
			method:       http.MethodPost,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
		},
		{
			name: "test follow user error",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceFollowUserErrorMock{},
			method:       http.MethodPost,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
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

			uh := handlers.NewUserHandler(tt.l, tt.userService)

			r := chi.NewRouter()

			r.Use(tt.middleware)

			r.Post("/follow-user", uh.FollowUser)

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

func TestUserHandler_UnfollowUser(t *testing.T) {
	// responses
	userNotFoundResponse, _ := json.Marshal(
		&handlers.SimpleResponse{Message: "User not found"},
	)

	internalError := handlers.MakeInternalErrorResponse()
	internalErrorResponse, _ := json.Marshal(&internalError)

	// tests cases
	tests := []struct {
		name         string
		l            *log.Logger
		userService  services.UserServiceI
		method       string
		statusCode   int
		responseBody string
		middleware   func(http.Handler) http.Handler
	}{
		{
			name: "test unfollow user success",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceUnfollowUserMock{},
			method:       http.MethodDelete,
			statusCode:   http.StatusNoContent,
			responseBody: "",
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
		},
		{
			name: "test user not passed in context",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceUnfollowUserMock{},
			method:       http.MethodDelete,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.IdentityMiddleware,
		},
		{
			name: "test user to unfollow not found",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceUnfollowUserNotFoundMock{},
			method:       http.MethodDelete,
			statusCode:   http.StatusNotFound,
			responseBody: string(userNotFoundResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
		},
		{
			name: "test get user to unfollow error",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceUnfollowUserGetUserErrorMock{},
			method:       http.MethodDelete,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
		},
		{
			name: "test unfollow user error",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceUnfollowUserErrorMock{},
			method:       http.MethodDelete,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
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

			uh := handlers.NewUserHandler(tt.l, tt.userService)

			r := chi.NewRouter()

			r.Use(tt.middleware)

			r.Delete("/unfollow-user", uh.UnfollowUser)

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

func TestUserHandler_GetCommunities(t *testing.T) {
	// responses
	communitiesResponse, _ := json.Marshal(
		&handlers.CommunitiesResponse{Communities: test_utils.CommunitiesList},
	)

	internalError := handlers.MakeInternalErrorResponse()
	internalErrorResponse, _ := json.Marshal(&internalError)

	// tests cases
	tests := []struct {
		name         string
		l            *log.Logger
		userService  services.UserServiceI
		method       string
		statusCode   int
		responseBody string
		middleware   func(http.Handler) http.Handler
	}{
		{
			name: "test get communities success",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceGetCommunitiesMock{},
			method:       http.MethodGet,
			statusCode:   http.StatusOK,
			responseBody: string(communitiesResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
		},
		{
			name: "test user not passed in context",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceGetCommunitiesMock{},
			method:       http.MethodGet,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.IdentityMiddleware,
		},
		{
			name: "test get communities errors",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceGetCommunitiesErrorMock{},
			method:       http.MethodGet,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				tt.method,
				"/get-communities",
				nil,
			)

			rr := httptest.NewRecorder()

			uh := handlers.NewUserHandler(tt.l, tt.userService)

			r := chi.NewRouter()

			r.Use(tt.middleware)

			r.Get("/get-communities", uh.GetCommunities)

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

func TestUserHandler_GetFollowers(t *testing.T) {
	// responses
	followersResponse, _ := json.Marshal(
		&handlers.FollowersResponse{Followers: test_utils.FollowersList},
	)

	internalError := handlers.MakeInternalErrorResponse()
	internalErrorResponse, _ := json.Marshal(&internalError)

	// tests cases
	tests := []struct {
		name         string
		l            *log.Logger
		userService  services.UserServiceI
		method       string
		statusCode   int
		responseBody string
		middleware   func(http.Handler) http.Handler
	}{
		{
			name: "test get followers success",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceGetFollowersMock{},
			method:       http.MethodGet,
			statusCode:   http.StatusOK,
			responseBody: string(followersResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
		},
		{
			name: "test user not passed in context",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceGetFollowersMock{},
			method:       http.MethodGet,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.IdentityMiddleware,
		},
		{
			name: "test get user errors",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceGetFollowersErrorMock{},
			method:       http.MethodGet,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
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

			uh := handlers.NewUserHandler(tt.l, tt.userService)

			r := chi.NewRouter()

			r.Use(tt.middleware)

			r.Get("/get-followers", uh.GetFollowers)

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
