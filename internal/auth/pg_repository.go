package auth

import (
	"context"

	"github.com/fekuna/go-rest-clean-architecture/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	FindByEmail(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	CreateAvatar(ctx context.Context, avatar *models.Avatar) (*models.Avatar, error)
	FindAvatarByID(ctx context.Context, avatarID uuid.UUID) (*models.Avatar, error)
	FindAvatarByFilePath(ctx context.Context, filePath string) (*models.Avatar, error)
}
