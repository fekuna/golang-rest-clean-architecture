package usecase

import (
	"context"
	"net/http"

	"github.com/fekuna/api-mc/config"
	"github.com/fekuna/api-mc/internal/auth"
	"github.com/fekuna/api-mc/internal/models"
	"github.com/fekuna/api-mc/pkg/httpErrors"
	"github.com/fekuna/api-mc/pkg/logger"
	"github.com/fekuna/api-mc/pkg/utils"
	"github.com/pkg/errors"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

// Auth UseCase
type authUC struct {
	cfg      *config.Config
	authRepo auth.Repository
	logger   logger.Logger
}

// Auth UseCase constructor
func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, log logger.Logger) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, logger: log}
}

// Create new user
func (u *authUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	// TODO: tracing

	existsUser, err := u.authRepo.FindByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrEmailAlreadyExists, nil)
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "auth.Register.GenerateJWTTOken"))
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Login user, returns user model with jwt token
func (u *authUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	// TODO: tracing

	foundUser, err := u.authRepo.FindByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePassword(user.Password); err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "authUC.GetUsers.ComparePasswords"))
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "authUC.GetUsers.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

// Find users by name
func (u *authUC) FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	// TODO: tracing

	return u.authRepo.FindByName(ctx, name, query)
}
