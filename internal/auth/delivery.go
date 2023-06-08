package auth

import "github.com/labstack/echo/v4"

// Auth HTTP Handlers interface
type Handlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	GetMe() echo.HandlerFunc
	UploadAvatar() echo.HandlerFunc
	GetAvatar() echo.HandlerFunc
}
