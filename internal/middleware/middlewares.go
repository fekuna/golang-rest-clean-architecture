package middleware

import (
	"github.com/fekuna/api-mc/config"
	"github.com/fekuna/api-mc/internal/auth"
	"github.com/fekuna/api-mc/internal/session"
	"github.com/fekuna/api-mc/pkg/logger"
)

// Middleware manager
type MiddlewareManager struct {
	cfg     *config.Config
	logger  logger.Logger
	origins []string
	sessUC  session.UCSession
	authUC  auth.UseCase
}

func NewMiddlewareManager(sessUC session.UCSession, authUC auth.UseCase, cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{sessUC: sessUC, authUC: authUC, cfg: cfg, origins: origins, logger: logger}
}
