package errorHandler

import "github.com/labstack/echo/v4"

type ErrorDetails struct {
	StatusCode int         `json:"statusCode"`
	Message    interface{} `json:"message"`
}

type SuccessDetails struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
}

func Panic(statusCode int, message string) {
	panic(ErrorDetails{
		StatusCode: statusCode,
		Message:    message,
	})
}

func Success(c echo.Context, statusCode int, data interface{}, message string) (err error) {
	successDetails := SuccessDetails{
		Data:       data,
		StatusCode: statusCode,
		Message:    message,
	}
	err = c.JSON(statusCode, successDetails)
	return
}

func Error(c echo.Context, statusCode int, error interface{}) (err error) {
	errorDetails := ErrorDetails{
		StatusCode: statusCode,
		Message:    error,
	}
	err = c.JSON(statusCode, errorDetails)
	return
}
