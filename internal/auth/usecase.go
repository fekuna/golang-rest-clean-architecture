package auth

import (
	"context"

	"github.com/fekuna/api-mc/internal/models"
	"github.com/fekuna/api-mc/pkg/utils"
	"github.com/google/uuid"
)

type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
