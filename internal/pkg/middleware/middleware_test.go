package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kidkrub/assessment-tax/internal/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestBasicAuthenticate(t *testing.T) {
	credential := config.New().BasicCredential()
	testCases := []struct {
		auth struct {
			username string
			password string
		}
		wantStatusCode int
	}{
		{struct {
			username string
			password string
		}{credential.Username, credential.Password}, http.StatusOK},
		{struct {
			username string
			password string
		}{"user", "password"}, http.StatusUnauthorized},
	}

	for _, tc := range testCases {
		e := echo.New()
		mw := BasicAuthenticate()
		e.Use(middleware.BasicAuth(mw))
		e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "[]") })
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth(tc.auth.username, tc.auth.password)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, tc.wantStatusCode, rec.Code)
	}
}
