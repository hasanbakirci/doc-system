package helpers

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomValidator struct {
	validator *validator.Validate
}

func Validate(ctx echo.Context, request interface{}) (result interface{}, err error) {
	if err = ctx.Bind(request); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	v := validator.New()
	if err := v.Struct(request); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		description := "Validation Error: "
		for _, e := range err.(validator.ValidationErrors) {
			description += e.Field() + ": " + e.ActualTag() + "  "
		}
		return nil, echo.NewHTTPError(http.StatusBadRequest, description)
	}
	return request, nil
}
