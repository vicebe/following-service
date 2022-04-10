package middleware_test

import (
	"github.com/go-chi/chi/v5"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/middleware"
	"github.com/vicebe/following-service/middleware/test_utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserMiddleware(t *testing.T) {
	r := chi.NewRouter()

	finalHandler := func(writer http.ResponseWriter, request *http.Request) {
		_, ok := request.Context().Value("user").(*data.User)

		if !ok {
			t.Fatal("User not found in context")
		}

		writer.WriteHeader(http.StatusNoContent)
	}

	r.Route("/{userID}", func(r chi.Router) {

		r.Route("/found-user", func(r chi.Router) {
			r.Use(middleware.GetUserMiddleware(test_utils.FoundUserMock{}))
			r.Get("/", finalHandler)
		})

		r.Route("/not-found-user", func(r chi.Router) {
			r.Use(middleware.GetUserMiddleware(test_utils.NotFoundUserMock{}))
			r.Get("/", finalHandler)
		})

		r.Route("/internal-error", func(r chi.Router) {
			r.Use(middleware.GetUserMiddleware(test_utils.UserServiceErrorMock{}))
			r.Get("/", finalHandler)
		})
	})

	t.Run("user found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/1/found-user", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusNoContent {
			t.Fatal(
				"Status not expected, expected: ",
				http.StatusNoContent,
				" got: ",
				res.StatusCode,
			)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/USER-NOT-FOUND/not-found-user",
			nil,
		)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusNotFound {
			t.Fatal(
				"Status not expected, expected: ",
				http.StatusNotFound,
				" got: ",
				res.StatusCode,
			)
		}

	})

	t.Run("internal server error", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/1/internal-error",
			nil,
		)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal(
				"Status not expected, expected: ",
				http.StatusInternalServerError,
				" got: ",
				res.StatusCode,
			)
		}

	})

}
