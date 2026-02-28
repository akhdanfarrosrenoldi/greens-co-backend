package handlers

import (
	"math"
	"net/http"
	"strconv"

	"greens-co/backend/internal/dto"
	"greens-co/backend/internal/repositories"
	"greens-co/backend/internal/services"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	svc *services.ProductService
}

func NewProductHandler(svc *services.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) GetAll(c echo.Context) error {
	available := true
	availStr := c.QueryParam("available")
	if availStr == "false" {
		available = false
	}

	min, _ := strconv.ParseInt(c.QueryParam("min"), 10, 64)
	max, _ := strconv.ParseInt(c.QueryParam("max"), 10, 64)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 9
	}

	f := repositories.ProductFilter{
		Cat:       c.QueryParam("cat"),
		Sort:      c.QueryParam("sort"),
		Min:       min,
		Max:       max,
		Available: &available,
		Search:    c.QueryParam("search"),
		Page:      page,
		Limit:     limit,
	}

	if availStr == "" {
		f.Available = &available // default true
	}

	products, total, err := h.svc.GetAll(f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

func (h *ProductHandler) GetBySlug(c echo.Context) error {
	product, err := h.svc.GetBySlug(c.Param("slug"))
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Product not found"})
	}
	return c.JSON(http.StatusOK, dto.Response{Data: product})
}
