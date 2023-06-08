package session

import (
	"context"

	"github.com/fekuna/go-rest-clean-architecture/internal/models"
)

type UseCase interface {
	CreateSession(ctx context.Context, session *models.Session) (*models.Session, error)
	UpsertSession(ctx context.Context, session *models.Session) (*models.Session, error)
}
