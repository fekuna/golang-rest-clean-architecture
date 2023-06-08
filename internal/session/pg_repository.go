package session

import (
	"context"

	"github.com/fekuna/go-rest-clean-architecture/internal/models"
)

// Session repository
type Repository interface {
	CreateSession(ctx context.Context, sess *models.Session) (*models.Session, error)
	UpdateSessionByUserId(ctx context.Context, sess *models.Session) (*models.Session, error)
	FindSessionByUserId(ctx context.Context, sess *models.Session) (*models.Session, error)
}
