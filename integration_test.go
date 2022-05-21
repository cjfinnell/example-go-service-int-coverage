//nolint:paralleltest,lll
package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	is := is.New(t)

	srv := newServer(&config{RedisURL: "redis:6379"})

	for _, step := range []struct {
		name           string
		method         string
		path           string
		auth           bool
		expectedStatus int
		expectedResp   string
	}{
		{name: "no auth get", method: http.MethodGet, path: "/testKey", expectedStatus: http.StatusUnauthorized},
		{name: "no auth set", method: http.MethodPost, path: "/testKey/testValue", expectedStatus: http.StatusUnauthorized},
		{name: "no auth del", method: http.MethodDelete, path: "/testKey", expectedStatus: http.StatusUnauthorized},
		{name: "get before set", method: http.MethodGet, path: "/testKey", expectedStatus: http.StatusNotFound, auth: true},
		{name: "set", method: http.MethodPost, path: "/testKey/testValue", expectedStatus: http.StatusCreated, auth: true},
		{name: "get", method: http.MethodGet, path: "/testKey", expectedStatus: http.StatusOK, expectedResp: "testValue", auth: true},
		{name: "del", method: http.MethodDelete, path: "/testKey", expectedStatus: http.StatusOK, auth: true},
		{name: "get after del", method: http.MethodGet, path: "/testKey", expectedStatus: http.StatusNotFound, auth: true},
	} {
		t.Run(step.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(step.method, step.path, nil)
			if step.auth {
				req.Header.Set("Authorization", "test-auth")
			}
			srv.router.ServeHTTP(w, req)
			is.Equal(step.expectedStatus, w.Code)
			if step.expectedResp != "" {
				is.Equal(step.expectedResp, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}
