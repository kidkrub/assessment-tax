package middleware

import (
	"github.com/kidkrub/assessment-tax/internal/pkg/config"
	"github.com/labstack/echo/v4"
)

func BasicAuthenticate() func(username, password string, c echo.Context) (bool, error) {
	return func(username string, password string, c echo.Context) (bool, error) {
		credential := config.New().BasicCredential()
		if username == credential.Username && password == credential.Password {
			return true, nil
		}
		return false, nil
	}
}
