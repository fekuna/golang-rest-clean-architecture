// go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package auth

import (
	"context"

	"github.com/fekuna/api-mc/internal/models"
)

// Auth repository interface
type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	FindByEmail(ctx context.Context, user *models.User) (*models.User, error)
}
