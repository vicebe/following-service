package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func getAccessTokenFromKeycloak(iamURL string) (string, error) {
	body := url.Values{}
	body.Add("username", "admin")
	body.Add("password", "admin")
	body.Add("grant_type", "password")
	body.Add("client_id", "admin-cli")
	resp, err := http.PostForm(iamURL, body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Request not ok when trying to get access token")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)

	if err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)

	if !ok {
		return "", fmt.Errorf("Could not get access token from response.")
	}

	return accessToken, nil
}

func TestAuthorizationMiddleware(t *testing.T) {

	iamCertURL := "http://localhost:8080/realms/master/protocol/openid-connect/certs"
	iamAccessTokenURL := "http://localhost:8080/realms/master/protocol/openid-connect/token"

	// starting docker IAM service using keycloak
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "Could not connect to docker")

	opts := &dockertest.RunOptions{
		Hostname:     "localhost",
		Repository:   "quay.io/keycloak/keycloak",
		Tag:          "19.0.2",
		Cmd:          []string{"start-dev"},
		ExposedPorts: []string{"8080"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8080/tcp": {{HostIP: "127.0.0.1", HostPort: "8080/tcp"}},
		},
		Env: []string{"KEYCLOAK_ADMIN_PASSWORD=admin", "KEYCLOAK_ADMIN=admin"},
	}
	resource, err := pool.RunWithOptions(opts)
	require.NoError(t, err, "could not start container")

	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	})

	var resp *http.Response

	pool.MaxWait = 5 * time.Minute
	err = pool.Retry(func() error {
		resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("8080/tcp"), "/"))
		if err != nil {
			t.Log("container not ready, waiting...")
			return err
		}
		return nil
	})
	require.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	// Creating router and server
	r := chi.NewRouter()

	r.With(AuthorizationMiddleware(iamCertURL)).Get("/protected-route", func(writer http.ResponseWriter, request *http.Request) {
		token, ok := request.Context().Value("token").(jwt.Token)
		require.True(t, ok, "Error getting token")
		require.NotEmpty(t, token, "Empty token")

		writer.WriteHeader(http.StatusOK)
	})

	// tests start
	t.Run("valid token", func(t *testing.T) {
		token, err := getAccessTokenFromKeycloak(iamAccessTokenURL)
		require.NoError(t, err, "Error getting access token from server")
		t.Log(token)
		req := httptest.NewRequest(http.MethodGet, "/protected-route", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		respBodyRaw, err := io.ReadAll(res.Body)
		require.NoError(t, err, "could not read response body")
		respBody := string(respBodyRaw)
		t.Log(respBody)
		require.Equal(t, http.StatusOK, res.StatusCode, "status not ok")
	})

	t.Run("no token", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/protected-route", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		respBodyRaw, err := io.ReadAll(res.Body)
		require.NoError(t, err, "could not read response body")
		respBody := string(respBodyRaw)
		require.Equal(t, http.StatusForbidden, res.StatusCode, "status not forbidden")
		require.Equal(t, "{\"message\":\"Token not provided in request\"}\n", respBody)

		req = httptest.NewRequest(http.MethodGet, "/protected-route", nil)
		req.Header.Add("Authorization", "Bearer ")
		rr = httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res = rr.Result()

		respBodyRaw, err = io.ReadAll(res.Body)
		require.NoError(t, err, "could not read response body")
		respBody = string(respBodyRaw)
		require.Equal(t, http.StatusForbidden, res.StatusCode, "status not forbidden")
		require.Equal(t, "{\"message\":\"Token not provided in request\"}\n", respBody)
	})

	t.Run("invalid token", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/protected-route", nil)
		req.Header.Add("Authorization", "Bearer invalid-jwt-token")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		respBodyRaw, err := io.ReadAll(res.Body)
		require.NoError(t, err, "could not read response body")
		respBody := string(respBodyRaw)
		require.Equal(t, http.StatusForbidden, res.StatusCode, "status not forbidden")
		require.Equal(t, "{\"message\":\"not authorized\"}\n", respBody)

	})
}
