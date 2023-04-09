package http

import (
	"context"
	"net/http"

	"github.com/fekuna/api-mc/config"
	"github.com/fekuna/api-mc/internal/auth"
	"github.com/fekuna/api-mc/internal/models"
	"github.com/fekuna/api-mc/internal/session"
	"github.com/fekuna/api-mc/pkg/httpErrors"
	"github.com/fekuna/api-mc/pkg/logger"
	"github.com/fekuna/api-mc/pkg/utils"
	"github.com/labstack/echo/v4"
)

type authHandlers struct {
	cfg    *config.Config
	logger logger.Logger
	authUC auth.UseCase
	sessUC session.UCSession
}

// NewAuthHandlers Auth handlers constructor
func NewAuthHandlers(cfg *config.Config, log logger.Logger, authUC auth.UseCase, sessUC session.UCSession) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, logger: log, sessUC: sessUC}
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
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdUser, err := h.authUC.Register(ctx, user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		sess, err := h.sessUC.CreateSession(ctx, &models.Session{
			UserID: createdUser.User.UserID,
		}, h.cfg.Session.Expire)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		c.SetCookie(utils.CreateSessionCookie(h.cfg, sess))

		return c.JSON(http.StatusCreated, createdUser)
	}
}
