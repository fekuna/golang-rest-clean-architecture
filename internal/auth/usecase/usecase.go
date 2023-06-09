package usecase

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/fekuna/go-rest-clean-architecture/config"
	"github.com/fekuna/go-rest-clean-architecture/internal/auth"
	"github.com/fekuna/go-rest-clean-architecture/internal/models"
	"github.com/fekuna/go-rest-clean-architecture/pkg/httpErrors"
	"github.com/fekuna/go-rest-clean-architecture/pkg/logger"
	"github.com/fekuna/go-rest-clean-architecture/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Auth Usecase
type authUC struct {
	cfg       *config.Config
	logger    logger.Logger
	authRepo  auth.Repository
	minioRepo auth.MinioRepository
}

// Auth usecase constructor
func NewAuthUseCase(cfg *config.Config, logger logger.Logger, authRepo auth.Repository, minioRepo auth.MinioRepository) auth.UseCase {
	return &authUC{
		cfg:       cfg,
		logger:    logger,
		authRepo:  authRepo,
		minioRepo: minioRepo,
	}
}

func (u *authUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	// TODO: Tracing

	existsUser, err := u.authRepo.FindByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrEmailAlreadyExists, err)
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser.SanitizePassword()

	accessToken, err := utils.GenerateJWTToken(createdUser, u.cfg, time.Minute*30)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.AccessToken.GenerateJWTTOken"))
	}

	refreshToken, err := utils.GenerateJWTToken(createdUser, u.cfg, (time.Hour*24)*30)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.RefreshToken.GenerateJWTTOken"))
	}

	authToken := models.AuthToken{
		AccesToken:   accessToken,
		RefreshToken: refreshToken,
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: authToken,
	}, nil
}

func (u *authUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	// TODO: tracing

	foundUser, err := u.authRepo.FindByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePassword(user.Password); err != nil {
		fmt.Println("mashok")
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "authUC.GetUsers.ComparePassword"))
	}

	foundUser.SanitizePassword()

	accessToken, err := utils.GenerateJWTToken(foundUser, u.cfg, time.Minute*30)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.AccessToken.GenerateJWTToken"))
	}

	refreshToken, err := utils.GenerateJWTToken(foundUser, u.cfg, (time.Hour*24)*30)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.RefreshToken.GenerateJWTToken"))
	}

	authToken := models.AuthToken{
		AccesToken:   accessToken,
		RefreshToken: refreshToken,
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: authToken,
	}, nil
}

func (u *authUC) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	// TODO: tracing

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.SanitizePassword()

	return user, nil
}

// Upload user avatar
func (u *authUC) UploadAvatar(ctx context.Context, userID uuid.UUID, file models.UploadInput) (*models.User, error) {
	// TODO: Tracing

	uploadInfo, err := u.minioRepo.PutObject(ctx, file)
	if err != nil {
		u.logger.Errorf("AuthUC.Update.UploadAvatar: %s", err)
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.UploadAvatar.PutObject"))
	}

	avatarURL := u.generateMinioURL(file.BucketName, uploadInfo.Key)

	updatedUser, err := u.authRepo.Update(ctx, &models.User{
		UserID: userID,
		Avatar: &avatarURL,
	})
	if err != nil {
		return nil, err
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}

func (u *authUC) GetAvatar(ctx context.Context) (*url.URL, error) {
	return u.minioRepo.GetObjectUrl(ctx, "static", "static/pathnich/6b2dba16-c2eb-4cce-89c1-7d2b4c5368fa-camera_lense_0.jpeg", time.Hour*24*7)
}

func (u *authUC) generateMinioURL(bucket string, key string) string {
	return fmt.Sprintf("%s/%s/%s", u.cfg.Minio.Endpoint, bucket, key)
}
