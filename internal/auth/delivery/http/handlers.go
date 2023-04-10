package http

import (
	"context"
	"errors"
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

// Login godoc
// @Summary Login new user
// @Description login user, returns user and set session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/login [post]
func (h *authHandlers) Login() echo.HandlerFunc {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}
	return func(c echo.Context) error {
		// TODO: tracing
		ctx := context.Background()

		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userWithToken, err := h.authUC.Login(ctx, &models.User{
			Email:    login.Email,
			Password: login.Password,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		sess, err := h.sessUC.CreateSession(ctx, &models.Session{
			UserID: userWithToken.User.UserID,
		}, h.cfg.Session.Expire)

		c.SetCookie(utils.CreateSessionCookie(h.cfg, sess))

		return c.JSON(http.StatusOK, userWithToken)
	}
}

// Logout godoc
// @Summary Logout user
// @Description logout user removing session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /auth/logout [post]
func (h *authHandlers) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: tracing

		ctx := context.Background()

		cookie, err := c.Cookie("session-id")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(err))
		}

		if err := h.sessUC.DeleteByID(ctx, cookie.Value); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		utils.DeleteSessionCookie(c, h.cfg.Session.Name)

		return c.NoContent(http.StatusOK)
	}
}

// FindByName godoc
// @Summary Find by name
// @Description Find user by name
// @Tags Auth
// @Accept json
// @Param name query string false "username" Format(username)
// @Produce json
// @Success 200 {object} models.UsersList
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/find [get]
func (h *authHandlers) FindByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: tracing
		ctx := context.Background()

		if c.QueryParam("name") == "" {
			utils.LogResponseError(c, h.logger, httpErrors.NewBadRequestError("name is required"))
			return c.JSON(http.StatusBadRequest, httpErrors.NewBadRequestError("name is required"))
		}

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		response, err := h.authUC.FindByName(ctx, c.QueryParam("name"), paginationQuery)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, response)
	}
}
