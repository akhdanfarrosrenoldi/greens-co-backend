package handlers

import (
	"net/http"

	"greens-co/backend/internal/dto"
	"greens-co/backend/internal/services"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	svc *services.CategoryService
}

func NewCategoryHandler(svc *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) GetAll(c echo.Context) error {
	categories, err := h.svc.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Data: categories})
}
