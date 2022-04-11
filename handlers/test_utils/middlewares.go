package test_utils

import (
	"context"
	"github.com/vicebe/following-service/data"
	"net/http"
)

func IdentityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			next.ServeHTTP(writer, request)
		})
}

func AddUserToRequestContext(user *data.User) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				ctx := context.WithValue(
					request.Context(),
					"user",
					user,
				)
				next.ServeHTTP(writer, request.WithContext(ctx))
			})
	}
}
func AddCommunityToRequestContext(
	community *data.Community,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				ctx := context.WithValue(
					request.Context(),
					"community",
					community,
				)
				next.ServeHTTP(writer, request.WithContext(ctx))
			})
	}
}
