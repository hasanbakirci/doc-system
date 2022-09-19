package auth

import (
	"net/http"

	"github.com/hasanbakirci/doc-system/pkg/errorHandler"
	"github.com/hasanbakirci/doc-system/pkg/helpers"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func (h *Handler) createUser(c echo.Context) error {
	request := new(CreateUserRequest)
	if _, err := helpers.Validate(c, request); err != nil {
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	//if err := c.Bind(request); err != nil {
	//	return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	//}
	result, err := h.service.Create(c.Request().Context(), *request)
	if err != nil {
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	return errorHandler.Success(c, http.StatusCreated, result, "Success")
}

func (h *Handler) updateUser(c echo.Context) error {
	id := c.Param("id")
	request := new(UpdateUserRequest)
	if _, err := helpers.Validate(c, request); err != nil {
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	//if err := c.Bind(request); err != nil {
	//	return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	//}

	result, err := h.service.Update(c.Request().Context(), id, *request)
	if !result {
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	return errorHandler.Success(c, http.StatusOK, result, "Success")
}

func (h *Handler) deleteUser(c echo.Context) error {
	id := c.Param("id")

	result, err := h.service.Delete(c.Request().Context(), id)
	if !result {
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	return errorHandler.Success(c, http.StatusOK, result, "Success")
}

func (h *Handler) getAllUsers(c echo.Context) error {
	result, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return errorHandler.Error(c, http.StatusNotFound, err.Error())
	}
	return errorHandler.Success(c, http.StatusOK, result, "Success")
}

func (h *Handler) getByIdUser(c echo.Context) error {
	id := c.Param("id")

	result, err := h.service.GetById(c.Request().Context(), id)
	if err != nil {
		return errorHandler.Error(c, http.StatusNotFound, err.Error())
	}
	return errorHandler.Success(c, http.StatusOK, result, "Success")
}

func (h *Handler) loginUser(c echo.Context) error {
	request := new(LoginUserRequest)
	if err := c.Bind(request); err != nil {
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	result, err := h.service.Login(c.Request().Context(), *request)
	if err != nil {
		return errorHandler.Error(c, http.StatusNotFound, err.Error())
	}
	return errorHandler.Success(c, http.StatusOK, result, "Success")
}

func NewUserHandler(s Service) Handler {
	return Handler{service: s}
}

func RegisterUserHandlers(instance *echo.Echo, h Handler) {
	instance.POST("api/users", h.createUser)
	instance.PUT("api/users/:id", h.updateUser)
	instance.DELETE("api/users/:id", h.deleteUser)
	instance.GET("api/users", h.getAllUsers)
	instance.GET("api/users/:id", h.getByIdUser)
	instance.POST("api/users/login", h.loginUser)
}
