package http

import (
	"github.com/fekuna/go-rest-clean-architecture/internal/auth"
	"github.com/fekuna/go-rest-clean-architecture/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.Use(mw.AuthJWTMiddleware)
	authGroup.GET("/me", h.GetMe())
	authGroup.POST("/:user_id/avatar", h.UploadAvatar(), mw.UploadFileMiddleware)
	authGroup.GET("/avatar", h.GetAvatar())
}
