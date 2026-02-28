package handlers

import (
	"net/http"

	"greens-co/backend/internal/dto"
	"greens-co/backend/internal/services"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	svc *services.OrderService
}

func NewOrderHandler(svc *services.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) Create(c echo.Context) error {
	var req dto.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}

	order, err := h.svc.Create(c, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, dto.Response{Data: order, Message: "Order created"})
}

func (h *OrderHandler) GetMyOrders(c echo.Context) error {
	orders, err := h.svc.GetMyOrders(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Data: orders})
}
