package auth

import (
	"context"
	"net/url"

	"github.com/fekuna/go-rest-clean-architecture/internal/models"
	"github.com/google/uuid"
)

type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UploadAvatar(ctx context.Context, userID uuid.UUID, file models.UploadInput) (*models.User, error)
	GetAvatarURL(ctx context.Context, avatarID uuid.UUID) (*url.URL, error)
}
