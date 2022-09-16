package middleware

import (
	"github.com/hasanbakirci/doc-system/pkg/helpers"
	"github.com/hasanbakirci/doc-system/pkg/response"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func TokenHandlerMiddlewareFunc(secret string, roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("Authorization") != "" {

				token := strings.Split(c.Request().Header.Get("Authorization"), " ")
				if token[0] != "Bearer" {
					log.Error("Authorization type is not Bearer.")
					return response.Error(c, http.StatusUnauthorized, "Authorization type is not Bearer.")
				}
				claims := helpers.VerifyToken(token[1], secret)
				if claims == nil {
					log.Error("The token is not correct.")
					return response.Error(c, http.StatusUnauthorized, "The token is not correct.")
				}

				if !handleRoles(roles, claims.Role) {
					log.Error("The user's role is not equal to the expected role.")
					return response.Error(c, http.StatusForbidden, "The user's role is not equal to the expected role.")
				}

				c.Set("id", claims.ID)
				log.Info("id field in context is set to : %s", claims.ID)
				return next(c)
			}
			log.Error("Authorization header is empty.")
			return response.Error(c, http.StatusUnauthorized, "Authorization header is empty.")
		}
	}
}

func handleRoles(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
