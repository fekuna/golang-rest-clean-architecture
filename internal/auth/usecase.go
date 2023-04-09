package auth

import (
	"context"

	"github.com/fekuna/api-mc/internal/models"
)

type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
}
