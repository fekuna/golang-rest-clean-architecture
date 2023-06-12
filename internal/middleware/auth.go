package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/fekuna/go-rest-clean-architecture/pkg/httpErrors"
	"github.com/fekuna/go-rest-clean-architecture/pkg/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// JWT way of auth using cookie or Authorization header
func (mw *MiddlewareManager) AuthJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

			if err := mw.validateJWTToken(c, tokenString); err != nil {
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

		if err = mw.validateJWTToken(c, cookie.Value); err != nil {
			mw.logger.Errorf("validateJWTToken", err.Error())
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}
		return next(c)
	}
}

func (mw *MiddlewareManager) validateJWTToken(c echo.Context, tokenString string) error {
	if tokenString == "" {
		return httpErrors.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(mw.cfg.Server.JwtSecretKey)
		return secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return httpErrors.InvalidJWTClaims
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

		u, err := mw.authUC.GetByID(c.Request().Context(), userUUID)
		if err != nil {
			return err
		}

		c.Set("user", u)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, u)
		// req := c.Request().WithContext(ctx)
		c.SetRequest(c.Request().WithContext(ctx))
	}

	return nil
}

// for Docs

// func (mw *MiddlewareManager) AuthJWTMiddleware(authUC auth.UseCase, cfg *config.Config) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			bearerHeader := c.Request().Header.Get("Authorization")

// 			mw.logger.Infof("auth middleware bearerHeader %s", bearerHeader)

// 			if bearerHeader != "" {
// 				headerParts := strings.Split(bearerHeader, " ")
// 				if len(headerParts) != 2 {
// 					mw.logger.Error("auth middleware", zap.String("headerParts", "len(headerParts) != 2"))
// 					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
// 				}

// 				tokenString := headerParts[1]

// 				if err := mw.validateJWTToken(tokenString, authUC, c, cfg); err != nil {
// 					mw.logger.Error("middleware validateJWTToken", zap.String("headerJWT", err.Error()))
// 					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
// 				}

// 				return next(c)
// 			}

// 			cookie, err := c.Cookie("jwt-token")
// 			if err != nil {
// 				mw.logger.Errorf("c.Cookie", err.Error())
// 				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
// 			}

// 			if err = mw.validateJWTToken(cookie.Value, authUC, c, cfg); err != nil {
// 				mw.logger.Errorf("validateJWTToken", err.Error())
// 				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
// 			}
// 			return next(c)
// 		}
// 	}
// }
