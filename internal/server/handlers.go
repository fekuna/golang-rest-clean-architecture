package server

import (
	"net/http"

	authHttp "github.com/fekuna/go-rest-clean-architecture/internal/auth/delivery/http"
	authRepository "github.com/fekuna/go-rest-clean-architecture/internal/auth/repository"
	authUseCase "github.com/fekuna/go-rest-clean-architecture/internal/auth/usecase"
	apiMiddlewares "github.com/fekuna/go-rest-clean-architecture/internal/middleware"
	sessRepository "github.com/fekuna/go-rest-clean-architecture/internal/session/repository"
	"github.com/fekuna/go-rest-clean-architecture/internal/session/usecase"
	"github.com/fekuna/go-rest-clean-architecture/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Map Server Handlers
func (s *Server) MapHandlers(e *echo.Echo) error {
	// TODO: setup metrics

	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	sRepo := sessRepository.NewSessionRepository(s.redisClient, s.cfg)
	authRedisRepo := authRepository.NewAuthRedisRepo(s.redisClient)
	aAWSRepo := authRepository.NewAuthAWSRepository(s.awsClient)

	// Init useCase
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, aAWSRepo, s.logger)
	sessUC := usecase.NewSessionUseCase(sRepo, s.cfg)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, sessUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(sessUC, authUC, s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)

	if s.cfg.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))

	e.Use(middleware.RequestID())

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
