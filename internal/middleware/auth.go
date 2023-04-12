package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/fekuna/api-mc/config"
	"github.com/fekuna/api-mc/internal/auth"
	"github.com/fekuna/api-mc/pkg/httpErrors"
	"github.com/fekuna/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// JWT way of auth using cookie or Authorization header
func (mw *MiddlewareManager) AuthJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")

			mw.logger.Infof("auth middleware bearerHeader %s", bearerHeader)

			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					mw.logger.Error("auth middleware", zap.String("headerParts", "len(headerParts) != 2"))
					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
				}

				tokenString := headerParts[1]

				if err := mw.validateJWTToken(tokenString, mw.authUC, c, mw.cfg); err != nil {
					mw.logger.Error("middleware validateJWTToken", zap.String("headerJWT", err.Error()))
					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
				}

				return next(c)
			}

			cookie, err := c.Cookie("jwt-token")
			if err != nil {
				mw.logger.Errorf("c.Cookie", err.Error())
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			if err = mw.validateJWTToken(cookie.Value, mw.authUC, c, mw.cfg); err != nil {
				mw.logger.Errorf("validateJWTToken", err.Error())
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			return next(c)
		}
	}
}

// Auth sessions middleware using redis
func (mw *MiddlewareManager) AuthSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(mw.cfg.Session.Name)
		if err != nil {
			mw.logger.Errorf("AuthSessionMiddleware RequestID: %s, Error: %s",
				utils.GetRequestID(c),
				err.Error(),
			)
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
			}
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		sid := cookie.Value

		sess, err := mw.sessUC.GetSessionByID(c.Request().Context(), cookie.Value)
		if err != nil {
			mw.logger.Errorf("GetSessionID RequestID: %s, Error: %s",
				utils.GetRequestID(c),
				err.Error(),
			)
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		user, err := mw.authUC.GetByID(c.Request().Context(), sess.UserID)
		if err != nil {
			mw.logger.Errorf("GetByID RequestID: %s, Error: %s",
				utils.GetRequestID(c),
				err.Error(),
			)
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		c.Set("sid", sid)
		c.Set("uid", sess.SessionID)
		c.Set("user", user)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, user)
		c.SetRequest(c.Request().WithContext(ctx))

		mw.logger.Info(
			"SessionMiddleware, RequestID: %s, IP: %s, UserID: %s, CookieSessionID: %s",
			utils.GetRequestID(c),
			utils.GetIPAddress(c),
			user.UserID.String(),
			cookie.Value,
		)

		return next(c)
	}
}

func (mw *MiddlewareManager) validateJWTToken(tokenString string, authUC auth.UseCase, c echo.Context, cfg *config.Config) error {
	if tokenString == "" {
		return httpErrors.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		secret := []byte(cfg.Server.JwtSecretKey)
		return secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return httpErrors.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return httpErrors.InvalidJWTClaims
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		u, err := authUC.GetByID(c.Request().Context(), userUUID)
		if err != nil {
			return err
		}

		c.Set("user", u)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, u)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}
