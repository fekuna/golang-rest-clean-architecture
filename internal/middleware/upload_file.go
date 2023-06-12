package middleware

import (
	"bytes"
	"io"

	"github.com/fekuna/go-rest-clean-architecture/internal/models"
	"github.com/fekuna/go-rest-clean-architecture/pkg/httpErrors"
	"github.com/fekuna/go-rest-clean-architecture/pkg/utils"
	"github.com/labstack/echo/v4"
)

// JWT way of auth using cookie or Authorization header
func (mw *MiddlewareManager) UploadFileMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		image, err := utils.ReadImage(c, "file")
		if err != nil {
			utils.LogResponseError(c, mw.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		file, err := image.Open()
		if err != nil {
			utils.LogResponseError(c, mw.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		defer file.Close()

		binaryImage := bytes.NewBuffer(nil)
		if _, err = io.Copy(binaryImage, file); err != nil {
			utils.LogResponseError(c, mw.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		contentType, err := utils.CheckImageFileContentType(binaryImage.Bytes())
		if err != nil {
			utils.LogResponseError(c, mw.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		reader := bytes.NewReader(binaryImage.Bytes())

		uploadInput := models.UploadInput{
			File:        reader,
			Name:        image.Filename,
			Size:        image.Size,
			ContentType: contentType,
			BucketName:  mw.cfg.Minio.Bucket,
		}

		c.Set("uploadInput", uploadInput)

		return next(c)
	}
}
