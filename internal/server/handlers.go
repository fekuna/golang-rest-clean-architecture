package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	authHttp "github.com/fekuna/api-mc/internal/auth/delivery/http"
	authRepository "github.com/fekuna/api-mc/internal/auth/repository"
	authUseCase "github.com/fekuna/api-mc/internal/auth/usecase"
	apiMiddlewares "github.com/fekuna/api-mc/internal/middleware"
	sessionRepository "github.com/fekuna/api-mc/internal/session/repository"
	"github.com/fekuna/api-mc/internal/session/usecase"
	"github.com/fekuna/api-mc/pkg/csrf"
	"github.com/fekuna/api-mc/pkg/utils"
)

// Map Server Handlers
func (s *Server) MapHandlers(e *echo.Echo) error {
	// TODO: create metrics

	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	sRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)

	// Init useCase
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, s.logger)
	sessUC := usecase.NewSessionUseCase(sRepo, s.cfg)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, s.logger, authUC, sessUC)

	mw := apiMiddlewares.NewMiddlewareManager(sessUC, authUC, s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)

	// TODO: docs

	if s.cfg.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	// TODO: metrics middleware

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

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
