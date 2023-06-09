package httpResponse

import (
	"github.com/fekuna/go-rest-clean-architecture/pkg/httpErrors"
	"github.com/fekuna/go-rest-clean-architecture/pkg/logger"
	"github.com/fekuna/go-rest-clean-architecture/pkg/utils"
	"github.com/labstack/echo/v4"
)

type ApiResponse struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
}

func Success(c echo.Context, code int, data interface{}, message string) error {
	response := ApiResponse{
		Data:    data,
		Code:    code,
		Message: message,
		Success: true,
	}

	return c.JSON(code, response)
}

func Error(c echo.Context, err error) error {
	restErr := httpErrors.ParseError(err)

	response := ApiResponse{
		Data:    map[string]interface{}{},
		Code:    restErr.Status(),
		Message: restErr.Error(),
		Success: false,
	}

	return c.JSON(restErr.Status(), response)
}

func ErrorWithLog(c echo.Context, logger logger.Logger, err error) error {
	restErr := httpErrors.ParseError(err)

	response := ApiResponse{
		Data:    map[string]interface{}{},
		Code:    restErr.Status(),
		Message: restErr.Error(),
		Success: false,
	}

	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		utils.GetRequestID(c),
		utils.GetIPAddress(c),
		err,
	)

	return c.JSON(restErr.Status(), response)
}
