//nolint:lll
package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func FuzzIntegration(f *testing.F) {
	if testing.Short() {
		f.Skip("skipping integration test")
	}

	is := is.New(f)

	srv := newServer(&config{RedisURL: "redis:6379"})

	f.Add("key", "value")
	f.Add("k", "v")
	f.Add("joined-by-dash", "109234jx09234hvno8ynfcgy8n")

	f.Fuzz(func(t *testing.T, fuzzKey, fuzzValue string) {
		t.Logf("fuzz targets: '%s' '%s'", fuzzKey, fuzzValue)
		if strings.TrimSpace(fuzzKey) == "" {
			t.Skip()
		}
		if strings.TrimSpace(fuzzValue) == "" {
			t.Skip()
		}
		if url.PathEscape(fuzzKey+fuzzValue) != fuzzKey+fuzzValue {
			t.Skip()
		}

		routeFuzzKey := "/" + fuzzKey
		routeFuzzValue := routeFuzzKey + "/" + fuzzValue

		for _, step := range []struct {
			name           string
			method         string
			path           string
			auth           bool
			expectedStatus int
			expectedResp   string
		}{
			{name: "no auth get", method: http.MethodGet, path: routeFuzzKey, expectedStatus: http.StatusUnauthorized},
			{name: "no auth set", method: http.MethodPost, path: routeFuzzValue, expectedStatus: http.StatusUnauthorized},
			{name: "no auth del", method: http.MethodDelete, path: routeFuzzKey, expectedStatus: http.StatusUnauthorized},
			{name: "get before set", method: http.MethodGet, path: routeFuzzKey, expectedStatus: http.StatusNotFound, auth: true},
			{name: "set", method: http.MethodPost, path: routeFuzzValue, expectedStatus: http.StatusCreated, auth: true},
			{name: "get", method: http.MethodGet, path: routeFuzzKey, expectedStatus: http.StatusOK, expectedResp: fuzzValue, auth: true},
			{name: "del", method: http.MethodDelete, path: routeFuzzKey, expectedStatus: http.StatusOK, auth: true},
			{name: "get after del", method: http.MethodGet, path: routeFuzzKey, expectedStatus: http.StatusNotFound, auth: true},
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
	})
}
