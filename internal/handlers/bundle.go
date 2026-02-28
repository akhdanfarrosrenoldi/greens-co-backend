package handlers

import (
	"net/http"

	"greens-co/backend/internal/dto"
	"greens-co/backend/internal/services"

	"github.com/labstack/echo/v4"
)

type BundleHandler struct {
	svc *services.BundleService
}

func NewBundleHandler(svc *services.BundleService) *BundleHandler {
	return &BundleHandler{svc: svc}
}

func (h *BundleHandler) GetAll(c echo.Context) error {
	bundles, err := h.svc.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Data: bundles})
}
