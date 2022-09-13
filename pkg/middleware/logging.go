package middleware

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func LoggingMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := log.Fields{
			"method": c.Request().Method,
			"path":   c.Path(),
			"url":    c.Request().RequestURI,
		}
		log.WithFields(request).Info("request details")
		return next(c)
	}
}
