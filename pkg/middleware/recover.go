package middleware

import (
	"github.com/hasanbakirci/doc-system/pkg/errorHandler"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func RecoveryMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			//str := recover()
			//c.JSON(http.StatusInternalServerError, str)
			if r := recover(); r != nil {
				switch t := r.(type) {
				case errorHandler.ErrorDetails:
					log.Error(r)
					errorHandler.Error(c, t.StatusCode, t.Message)
				default:
					log.Warn(r)
					errorHandler.Error(c, 500, r)
				}
				//c.JSON(http.StatusInternalServerError, err)
				//errorHandler.ErrorResponse(c, 500, 5000, r)
			}
		}()
		return next(c)
	}
}
