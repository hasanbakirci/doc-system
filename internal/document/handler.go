package document

import (
	"github.com/hasanbakirci/doc-system/pkg/helpers"
	"github.com/hasanbakirci/doc-system/pkg/middleware"
	"github.com/hasanbakirci/doc-system/pkg/response"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Handler struct {
	service Service
}

func (h Handler) createDocument(c echo.Context) error {
	description := c.FormValue("description")
	file, err := c.FormFile("file")
	if err != nil {
		//return c.JSON(http.StatusBadRequest, err)
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	fileResult := helpers.AddFile(file)

	//if err := c.Bind(request); err != nil {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}
	result, err := h.service.Create(c.Request().Context(), CreateDocumentRequest{
		Name:        fileResult.FileName,
		Description: description,
		Extension:   fileResult.Extension,
		Path:        fileResult.Path,
		MimeType:    fileResult.MimeType,
	})

	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	//return c.JSON(http.StatusCreated, result)
	return response.Success(c, http.StatusCreated, result, "Success")
}

func (h Handler) updateDocument(c echo.Context) error {
	description := c.FormValue("description")
	value := c.Param("id")
	id, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	//request := new(UpdateDocumentRequest)
	//if err := c.Bind(request); err != nil {
	//	c.JSON(http.StatusBadRequest, err.Error())
	//}
	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	fileResult := helpers.AddFile(file)

	result, err := h.service.Update(c.Request().Context(), id, UpdateDocumentRequest{
		Name:        fileResult.FileName,
		Description: description,
		Extension:   fileResult.Extension,
		Path:        fileResult.Path,
		MimeType:    fileResult.MimeType,
	})
	if !result {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	return response.Success(c, http.StatusOK, result, "Success")
}

func (h Handler) deleteDocument(c echo.Context) error {
	value := c.Param("id")
	id, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	result, err := h.service.Delete(c.Request().Context(), id)
	if !result {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	return response.Success(c, http.StatusOK, result, "Success")
}

func (h Handler) getAllDocuments(c echo.Context) error {
	documents, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		//return c.JSON(http.StatusNotFound, err)
		return response.Error(c, http.StatusNotFound, err.Error())
	}
	return response.Success(c, http.StatusOK, documents, "Success")
}

func (h Handler) getByIdDocument(c echo.Context) error {
	value := c.Param("id")
	id, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	document, err := h.service.GetById(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	return response.Success(c, http.StatusOK, document, "Success")
}
func NewDocumentHandler(s Service) Handler {
	return Handler{service: s}
}

func RegisterDocumentHandlers(instance *echo.Echo, h Handler, secret string) {
	instance.POST("api/documents", h.createDocument)
	instance.PUT("api/documents/:id", h.updateDocument)
	instance.DELETE("api/documents/:id", h.deleteDocument)
	instance.GET("api/documents", h.getAllDocuments, middleware.TokenHandlerMiddlewareFunc(secret, "user", "admin"))
	instance.GET("api/documents/:id", h.getByIdDocument)
}
