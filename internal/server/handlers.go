package server

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	authHttp "github.com/fekuna/go-rest-clean-architecture/internal/auth/delivery/http"
	authRepository "github.com/fekuna/go-rest-clean-architecture/internal/auth/repository"
	authUC "github.com/fekuna/go-rest-clean-architecture/internal/auth/usecase"
	apiMiddlewares "github.com/fekuna/go-rest-clean-architecture/internal/middleware"
	sessRepository "github.com/fekuna/go-rest-clean-architecture/internal/session/repository"
	sessUC "github.com/fekuna/go-rest-clean-architecture/internal/session/usecase"
)

func (s *Server) MapHandlers(e *echo.Echo) error {

	// Init Repository
	authRepo := authRepository.NewAuthRepository(s.db)
	sessRepo := sessRepository.NewSessionRepository(s.db)
	// authMinioRepo := authRepository.NewAuthMinioRepository(s.minioConfig)

	// Init useCase
	authUC := authUC.NewAuthUseCase(s.cfg, s.logger, authRepo, *s.minioConfig)
	sessUC := sessUC.NewSessionUseCase(s.cfg, s.logger, sessRepo)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, s.logger, authUC, sessUC)

	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, s.logger, sessUC, authUC)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	if s.cfg.Server.Debug {
		e.Use(mw.DebugMiddleware)
	}

	v1 := e.Group("/api/v1")

	authGroup := v1.Group("/auth")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)

	return nil
}
