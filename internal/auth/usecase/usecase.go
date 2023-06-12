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
	"github.com/fekuna/go-rest-clean-architecture/pkg/db/minioS3"
	"github.com/fekuna/go-rest-clean-architecture/pkg/httpErrors"
	"github.com/fekuna/go-rest-clean-architecture/pkg/logger"
	"github.com/fekuna/go-rest-clean-architecture/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Auth Usecase
type authUC struct {
	cfg         *config.Config
	logger      logger.Logger
	authRepo    auth.Repository
	minioConfig minioS3.MinioConfig
}

// Auth usecase constructor
func NewAuthUseCase(cfg *config.Config, logger logger.Logger, authRepo auth.Repository, minioConfig minioS3.MinioConfig) auth.UseCase {
	return &authUC{
		cfg:         cfg,
		logger:      logger,
		authRepo:    authRepo,
		minioConfig: minioConfig,
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

	avatar, err := u.authRepo.FindAvatarByFilePath(ctx, "1b05aaaf-6881-4917-953e-048ac295a6e9-default.jpg")
	if err != nil {
		u.logger.Errorf("authUC.Register.FindAvatarByFilePath: %s", err)
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.FindAvatarByFilePath"))
	}

	user.AvatarID = avatar.AvatarID
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

	if err := u.minioConfig.CreateBucket(ctx, file.BucketName); err != nil {
		u.logger.Errorf("AuthUC.UploadAvatar.CreateBucket: %s", err)
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.UploadAvatar.CreateBucket"))
	}

	fileUploaded, err := u.minioConfig.PutObject(ctx, file)
	if err != nil {
		u.logger.Errorf("AuthUC.UploadAvatar.UploadAvatar: %s", err)
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.UploadAvatar.PutObject"))
	}

	// avatarURL := u.generateMinioURL(file.BucketName, uploadInfo.Key)

	createAvatar, err := u.authRepo.CreateAvatar(ctx, &models.Avatar{
		Bucket:   fileUploaded.Bucket,
		FilePath: fileUploaded.Key,
	})

	if err != nil {
		u.logger.Errorf("AuthUC.UploadAvatar.createAvatar: %s", err)
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.UploadAvatar.createAvatar"))
	}

	updatedUser, err := u.authRepo.Update(ctx, &models.User{
		UserID:   userID,
		AvatarID: createAvatar.AvatarID,
	})
	if err != nil {
		u.logger.Errorf("AuthUC.UploadAvatar.Update: %s", err)
		return nil, err
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}

func (u *authUC) GetAvatarURL(ctx context.Context, avatarID uuid.UUID) (*url.URL, error) {
	fmt.Println("avatarID: ", avatarID)
	avatar, err := u.authRepo.FindAvatarByID(ctx, avatarID)
	if err != nil {
		return nil, err
	}

	return u.minioConfig.GetObjectUrl(ctx, avatar.Bucket, avatar.FilePath, time.Hour*24*7)
}

func (u *authUC) generateMinioURL(bucket string, key string) string {
	return fmt.Sprintf("%s/%s/%s", u.cfg.Minio.Endpoint, bucket, key)
}
