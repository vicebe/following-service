package middleware

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/handlers"
	"net/http"
	"strings"
)

// obtainToken gets the authorization token from the bearer.
func obtainToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	tokenSplit := strings.Split(reqToken, "Bearer ")

	if len(tokenSplit) == 1 {
		return "", fmt.Errorf("token not provided")
	}

	token := tokenSplit[1]

	return token, nil
}

// AuthorizationMiddleware verifies if the token given is valid
func AuthorizationMiddleware(
	authCertsURL string,
) func(http.Handler) http.Handler {

	return func(handler http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := obtainToken(r)

			if err != nil {

				w.WriteHeader(http.StatusForbidden)

				response := handlers.SimpleResponse{
					Message: "Token not provided in request",
				}

				if err := data.ToJson(&response, w); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}

				return
			}

			set, err := jwk.Fetch(context.Background(), authCertsURL)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				response := handlers.SimpleResponse{
					Message: "Error Fetching Certs",
				}

				if err := data.ToJson(&response, w); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}

				return
			}

			jwtToken, err := jwt.Parse(
				[]byte(token),
				jwt.WithKeySet(set),
				jwt.WithValidate(true),
			)

			if err != nil {

				w.WriteHeader(http.StatusForbidden)

				response := handlers.SimpleResponse{
					Message: "not authorized",
				}

				if err := data.ToJson(&response, w); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}

				return
			}

			ctx := context.WithValue(r.Context(), "token", jwtToken)

			handler.ServeHTTP(w, r.WithContext(ctx))
		})

	}
}
