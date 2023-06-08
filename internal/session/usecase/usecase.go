package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fekuna/go-rest-clean-architecture/config"
	"github.com/fekuna/go-rest-clean-architecture/internal/models"
	"github.com/fekuna/go-rest-clean-architecture/internal/session"
	"github.com/fekuna/go-rest-clean-architecture/pkg/logger"
)

type SessionUC struct {
	cfg         *config.Config
	logger      logger.Logger
	sessionRepo session.Repository
}

func NewSessionUseCase(cfg *config.Config, logger logger.Logger, sessionRepo session.Repository) session.UseCase {
	return &SessionUC{
		cfg:         cfg,
		logger:      logger,
		sessionRepo: sessionRepo,
	}
}

func (s *SessionUC) CreateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	// TODO: Tracing

	return s.sessionRepo.CreateSession(ctx, session)
}

func (s *SessionUC) UpsertSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	// TODO: tracing

	_, err := s.sessionRepo.FindSessionByUserId(ctx, session)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// If No Data found. we insert it
	if errors.Is(err, sql.ErrNoRows) {
		return s.sessionRepo.CreateSession(ctx, session)
	} else {
		// update session if user has session
		return s.sessionRepo.UpdateSessionByUserId(ctx, session)
	}
}
