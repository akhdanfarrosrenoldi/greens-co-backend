package handlers

import (
	"net/http"

	"greens-co/backend/internal/dto"
	"greens-co/backend/internal/services"

	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	svc *services.PaymentService
}

func NewPaymentHandler(svc *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) Initiate(c echo.Context) error {
	var req dto.InitiatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}

	paymentURL, err := h.svc.Initiate(req.OrderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Data: dto.PaymentResponse{PaymentURL: paymentURL},
	})
}

func (h *PaymentHandler) Notification(c echo.Context) error {
	var notification map[string]interface{}
	if err := c.Bind(&notification); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid notification"})
	}

	if err := h.svc.HandleNotification(notification); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
}
