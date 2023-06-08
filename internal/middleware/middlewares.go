package middleware

import (
	"github.com/fekuna/go-rest-clean-architecture/config"
	"github.com/fekuna/go-rest-clean-architecture/internal/auth"
	"github.com/fekuna/go-rest-clean-architecture/internal/session"
	"github.com/fekuna/go-rest-clean-architecture/pkg/logger"
)

// Middleware manager
type MiddlewareManager struct {
	cfg    *config.Config
	logger logger.Logger
	sessUC session.UseCase
	authUC auth.UseCase
}

// Middleware manager constructor
func NewMiddlewareManager(
	cfg *config.Config,
	logger logger.Logger,
	sessUC session.UseCase,
	authUC auth.UseCase,
) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:    cfg,
		logger: logger,
		sessUC: sessUC,
		authUC: authUC,
	}
}
