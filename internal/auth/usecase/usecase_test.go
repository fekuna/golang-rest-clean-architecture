package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/fekuna/go-rest-clean-architecture/config"
	"github.com/fekuna/go-rest-clean-architecture/internal/auth/mock"
	"github.com/fekuna/go-rest-clean-architecture/internal/models"
	"github.com/fekuna/go-rest-clean-architecture/pkg/logger"
	"github.com/fekuna/go-rest-clean-architecture/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUC_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, nil, nil, apiLogger)

	user := &models.User{
		Email:    "email@gmail.com",
		Password: "123456",
	}

	ctx := context.Background()
	// TODO: Open Tracing

	mockAuthRepo.EXPECT().FindByEmail(ctx, gomock.Eq(user)).Return(nil, sql.ErrNoRows)
	mockAuthRepo.EXPECT().Register(ctx, gomock.Eq(user)).Return(user, nil)

	createdUser, err := authUC.Register(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.Nil(t, err)
}

func TestAuthUC_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, nil, apiLogger)

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", basePrefix, user.UserID)

	ctx := context.Background()
	// TODO: Open Tracing

	mockRedisRepo.EXPECT().GetByIDCtx(ctx, key).Return(nil, nil)
	mockAuthRepo.EXPECT().GetByID(ctx, gomock.Eq(user.UserID)).Return(user, nil)
	mockRedisRepo.EXPECT().SetUserCtx(ctx, key, cacheDuration, user).Return(nil)

	u, err := authUC.GetByID(ctx, user.UserID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, u)
}

func TestAuthUC_FindByName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, nil, apiLogger)

	userName := "name"
	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}
	ctx := context.Background()
	// TODO: Open Tracing

	usersList := &models.UsersList{}

	mockAuthRepo.EXPECT().FindByName(ctx, gomock.Eq(userName), query).Return(usersList, nil)

	userList, err := authUC.FindByName(ctx, userName, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, userList)
}

func TestAuthUC_GetUsers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, nil, apiLogger)

	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}
	ctx := context.Background()
	// TODO: Open Tracing

	usersList := &models.UsersList{}

	mockAuthRepo.EXPECT().GetUsers(ctx, query).Return(usersList, nil)

	users, err := authUC.GetUsers(ctx, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, users)
}

func TestAuthUC_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, nil, apiLogger)

	ctx := context.Background()
	// TODO: Open Tracing

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	mockUser := &models.User{
		Email:    "email@gmail.com",
		Password: string(hashPassword),
	}

	mockAuthRepo.EXPECT().FindByEmail(ctx, gomock.Eq(user)).Return(mockUser, nil)

	userWithToken, err := authUC.Login(ctx, user)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, userWithToken)
}
