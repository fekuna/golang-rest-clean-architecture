package http

import (
	"context"
	"net/http"
	"time"

	"github.com/fekuna/go-rest-clean-architecture/config"
	"github.com/fekuna/go-rest-clean-architecture/internal/auth"
	"github.com/fekuna/go-rest-clean-architecture/internal/models"
	"github.com/fekuna/go-rest-clean-architecture/internal/session"
	"github.com/fekuna/go-rest-clean-architecture/pkg/httpErrors"
	"github.com/fekuna/go-rest-clean-architecture/pkg/httpResponse"
	"github.com/fekuna/go-rest-clean-architecture/pkg/logger"
	"github.com/fekuna/go-rest-clean-architecture/pkg/utils"
	"github.com/labstack/echo/v4"
)

// Auth handlers
type authHandlers struct {
	cfg    *config.Config
	logger logger.Logger
	authUC auth.UseCase
	sessUC session.UseCase
}

func NewAuthHandlers(cfg *config.Config, logger logger.Logger, authUC auth.UseCase, sessUC session.UseCase) auth.Handlers {
	return &authHandlers{
		cfg:    cfg,
		logger: logger,
		authUC: authUC,
		sessUC: sessUC,
	}
}

// Register godoc
// @Summary Register new user
// @Description register new user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /auth/register [post]
func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: tracing
		ctx := context.Background()

		user := &models.User{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, h.logger, err)
			// return c.JSON(httpErrors.ErrorResponse(err))
			return httpResponse.Error(c, err)
		}

		createdUser, err := h.authUC.Register(ctx, user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			// return c.JSON(httpErrors.ErrorResponse(err))
			return httpResponse.Error(c, err)
		}

		session := &models.Session{
			RefreshToken: createdUser.Token.RefreshToken,
			ExpiresAt:    time.Now().Add(time.Hour * 24 * 30),
			UserID:       createdUser.User.UserID,
		}

		_, err = h.sessUC.CreateSession(ctx, session)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			// return c.JSON(httpErrors.ErrorResponse(err))
			return httpResponse.Error(c, err)
		}

		// return c.JSON(http.StatusCreated, createdUser)
		return httpResponse.Success(c, http.StatusCreated, createdUser.Token, "Success created user")
	}
}

// Login godoc
// @Summary Login user
// @Description login user, returns tokens and set session in DB
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/login [post]
func (h *authHandlers) Login() echo.HandlerFunc {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}
	return func(c echo.Context) error {
		// TODO: tracing
		ctx := context.Background()

		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			utils.LogResponseError(c, h.logger, err)
			// return c.JSON(httpErrors.ErrorResponse(err))
			return httpResponse.Error(c, err)
		}

		userWithToken, err := h.authUC.Login(ctx, &models.User{
			Email:    login.Email,
			Password: login.Password,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			// return c.JSON(httpErrors.ErrorResponse(err))
			return httpResponse.Error(c, err)
		}

		sess := &models.Session{
			RefreshToken: userWithToken.Token.RefreshToken,
			ExpiresAt:    time.Now().Add(time.Hour * 24 * 30),
			UserID:       userWithToken.User.UserID,
		}

		// TODO: can generate multiple session for multiple devices. for the future feature.
		_, err = h.sessUC.UpsertSession(ctx, sess)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			// return c.JSON(httpErrors.ErrorResponse(err))
			return httpResponse.Error(c, err)
		}

		// return c.JSON(http.StatusOK, userWithToken)
		return httpResponse.Success(c, http.StatusOK, userWithToken.Token, "")

	}
}

func (h *authHandlers) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: tracing
		ctx := context.Background()

		user, ok := c.Get("user").(*models.User)
		if !ok {
			// utils.LogResponseError(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			// return utils.ErrResponseWithLog(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			return httpResponse.ErrorWithLog(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		avatarUrl, err := h.authUC.GetAvatarURL(ctx, user.AvatarID)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}

		responseData := struct {
			AvatarUrl string `json:"avatar_url"`
			models.User
		}{
			AvatarUrl: avatarUrl.String(),
			User:      *user,
		}

		// return c.JSON(http.StatusOK, user)
		return httpResponse.Success(c, http.StatusOK, responseData, "Success Get User")
	}
}

// UploadAvatar godoc
// @Summary Post avatar
// @Description Post user avatar image
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param file formData file true "Body with image file"
// @Param bucket query string true "aws s3 bucket" Format(bucket)
// @Param id path int true "user_id"
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/{id}/avatar [post]
func (h *authHandlers) UploadAvatar() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: tracing

		ctx := context.Background()

		user, ok := c.Get("user").(*models.User)
		if !ok {
			// utils.LogResponseError(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			// return utils.ErrResponseWithLog(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			return httpResponse.ErrorWithLog(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		uploadInput := c.Get("uploadInput").(models.UploadInput)
		if !ok {
			return httpResponse.ErrorWithLog(c, h.logger, httpErrors.NewInternalServerError(httpErrors.InternalServerError))
		}

		updatedUser, err := h.authUC.UploadAvatar(ctx, user.UserID, uploadInput)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		// return c.JSON(http.StatusOK, updatedUser)
		return httpResponse.Success(c, http.StatusCreated, updatedUser, "avatar uploaded")
	}
}

// GetAvatar godoc
// @Summary Post avatar
// @Description Post user avatar image
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param file formData file true "Body with image file"
// @Param bucket query string true "aws s3 bucket" Format(bucket)
// @Param id path int true "user_id"
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/{id}/avatar [post]
func (h *authHandlers) GetAvatar() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: tracing

		user, ok := c.Get("user").(*models.User)
		if !ok {
			return httpResponse.ErrorWithLog(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		ctx := context.Background()

		url, err := h.authUC.GetAvatarURL(ctx, user.AvatarID)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}

		// return c.JSON(http.StatusOK, url.String())
		return httpResponse.Success(c, http.StatusOK, struct {
			AvatarUrl string `json:"avatar_url"`
		}{
			AvatarUrl: url.String(),
		}, "Success to get avatar")
	}
}
