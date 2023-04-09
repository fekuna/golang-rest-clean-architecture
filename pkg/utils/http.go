package utils

import (
	"fmt"
	"net/http"

	"github.com/fekuna/api-mc/config"
	"github.com/fekuna/api-mc/pkg/logger"
	"github.com/labstack/echo/v4"
)

// Get request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// Get user ip address
func GetIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}

// Get config path for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}

	return "./config/config-local"
}

// Read request body and validate
func ReadRequest(ctx echo.Context, request interface{}) error {
	fmt.Println("Before bind")
	if err := ctx.Bind(request); err != nil {
		fmt.Println("Error bind", err)
		return err
	}
	fmt.Println("After bind")
	return validate.StructCtx(ctx.Request().Context(), request)
}

// Error response with logging error for echo context
func LogResponseError(ctx echo.Context, logger logger.Logger, err error) {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		GetIPAddress(ctx),
		err,
	)
}

// Configure session cookie
func CreateSessionCookie(cfg *config.Config, session string) *http.Cookie {
	return &http.Cookie{
		Name:  cfg.Session.Name,
		Value: session,
		Path:  "/",
		// Domain: "/",
		// Expires: time.Now().Add(1 * time.Minute),
		RawExpires: "",
		MaxAge:     cfg.Session.Expire,
		Secure:     cfg.Cookie.Secure,
		HttpOnly:   cfg.Cookie.HTTPOnly,
		SameSite:   0,
	}
}
