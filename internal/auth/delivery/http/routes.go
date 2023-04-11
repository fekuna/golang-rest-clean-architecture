package http

import (
	"github.com/fekuna/api-mc/internal/auth"
	"github.com/fekuna/api-mc/internal/middleware"
	"github.com/labstack/echo/v4"
)

// Map auth routes
func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/logout", h.Logout())
	authGroup.GET("/find", h.FindByName())
	authGroup.Use(mw.AuthSessionMiddleware)
	authGroup.GET("/me", h.GetMe())
}
