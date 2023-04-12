// go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package auth

import (
	"context"

	"github.com/fekuna/api-mc/internal/models"
	"github.com/fekuna/api-mc/pkg/utils"
	"github.com/google/uuid"
)

// Auth repository interface
type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	FindByEmail(ctx context.Context, user *models.User) (*models.User, error)
	FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
}
