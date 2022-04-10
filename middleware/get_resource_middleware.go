package middleware

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/handlers"
	"github.com/vicebe/following-service/services"
	"net/http"
)

func GetUserMiddleware(
	service services.UserServiceI,
) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "userID")

			user, err := service.GetUser(userID)

			if err != nil {
				var httpStatus int
				var response handlers.SimpleResponse

				if err == data.ErrorUserNotFound {
					httpStatus = http.StatusNotFound
					response = handlers.SimpleResponse{
						Message: "User not found",
					}
				} else {
					httpStatus = http.StatusInternalServerError
					response = handlers.SimpleResponse{
						Message: "Something went wrong",
					}
				}

				rw.WriteHeader(httpStatus)

				if err := data.ToJson(&response, rw); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
				}

				return
			}

			ctx := context.WithValue(r.Context(), "user", user)

			next.ServeHTTP(rw, r.WithContext(ctx))

		})
	}
}

func GetCommunityMiddleware(
	service services.CommunityServiceI,
) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			communityID := chi.URLParam(r, "communityID")

			community, err := service.GetCommunity(communityID)

			if err != nil {
				var httpStatus int
				var response handlers.SimpleResponse

				if err == data.ErrorCommunityNotFound {
					httpStatus = http.StatusNotFound
					response = handlers.SimpleResponse{
						Message: "Community not found",
					}
				} else {
					httpStatus = http.StatusInternalServerError
					response = handlers.SimpleResponse{
						Message: "Something went wrong",
					}
				}

				rw.WriteHeader(httpStatus)

				if err := data.ToJson(&response, rw); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
				}

				return
			}

			ctx := context.WithValue(r.Context(), "community", community)

			next.ServeHTTP(rw, r.WithContext(ctx))

		})
	}
}
