package document

import (
	"fmt"
	"net/http"

	"github.com/hasanbakirci/doc-system/pkg/errorHandler"
	"github.com/hasanbakirci/doc-system/pkg/helpers"
	"github.com/hasanbakirci/doc-system/pkg/middleware"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func (h Handler) createDocument(c echo.Context) error {
	uid := c.Get("id")
	description := c.FormValue("description")
	file, err := c.FormFile("file")
	if err != nil {
		//return c.JSON(http.StatusBadRequest, err)
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	fileResult := helpers.AddFile(file, "text/plain; charset=utf-8")
	//if err := c.Bind(request); err != nil {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}
	result, err := h.service.Create(c.Request().Context(), CreateDocumentRequest{
		Name:        fileResult.FileName,
		Description: description,
		Extension:   fileResult.Extension,
		Path:        fileResult.Path,
		MimeType:    fileResult.MimeType,
	}, fmt.Sprintf("%v", uid))

	if err != nil {
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	//return c.JSON(http.StatusCreated, result)
	return errorHandler.Success(c, http.StatusCreated, result, "Success")
}

func (h Handler) updateDocument(c echo.Context) error {
	description := c.FormValue("description")
	id := c.Param("id")

	//request := new(UpdateDocumentRequest)
	//if err := c.Bind(request); err != nil {
	//	c.JSON(http.StatusBadRequest, err.Error())
	//}
	file, err := c.FormFile("file")
	if err != nil {
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	fileResult := helpers.AddFile(file, "text/plain; charset=utf-8")

	result, err := h.service.Update(c.Request().Context(), id, UpdateDocumentRequest{
		Name:        fileResult.FileName,
		Description: description,
		Extension:   fileResult.Extension,
		Path:        fileResult.Path,
		MimeType:    fileResult.MimeType,
	})
	if !result {
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	return errorHandler.Success(c, http.StatusOK, result, "Success")
}

func (h Handler) deleteDocument(c echo.Context) error {
	id := c.Param("id")

	result, err := h.service.Delete(c.Request().Context(), id)
	if !result {
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	return errorHandler.Success(c, http.StatusOK, result, "Success")
}

func (h Handler) getAllDocuments(c echo.Context) error {
	documents, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		//return c.JSON(http.StatusNotFound, err)
		return errorHandler.Error(c, http.StatusNotFound, err.Error())
	}
	return errorHandler.Success(c, http.StatusOK, documents, "Success")
}

func (h Handler) getByIdDocument(c echo.Context) error {
	id := c.Param("id")

	document, err := h.service.GetById(c.Request().Context(), id)
	if err != nil {
		return errorHandler.Error(c, http.StatusBadRequest, err.Error())
	}
	return errorHandler.Success(c, http.StatusOK, document, "Success")
}
func NewDocumentHandler(s Service) Handler {
	return Handler{service: s}
}

func RegisterDocumentHandlers(instance *echo.Echo, h Handler, secret string) {
	instance.POST("api/documents", h.createDocument, middleware.TokenHandlerMiddlewareFunc(secret, "user", "admin"))
	instance.PUT("api/documents/:id", h.updateDocument)
	instance.DELETE("api/documents/:id", h.deleteDocument)
	instance.GET("api/documents", h.getAllDocuments, middleware.TokenHandlerMiddlewareFunc(secret, "user", "admin"))
	instance.GET("api/documents/:id", h.getByIdDocument)
}
