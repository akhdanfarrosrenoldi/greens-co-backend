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

type AdminHandler struct {
	productSvc  *services.ProductService
	categorySvc *services.CategoryService
	bundleSvc   *services.BundleService
	orderSvc    *services.OrderService
	orderRepo   *repositories.OrderRepository
	productRepo *repositories.ProductRepository
}

func NewAdminHandler(
	productSvc *services.ProductService,
	categorySvc *services.CategoryService,
	bundleSvc *services.BundleService,
	orderSvc *services.OrderService,
	orderRepo *repositories.OrderRepository,
	productRepo *repositories.ProductRepository,
) *AdminHandler {
	return &AdminHandler{
		productSvc:  productSvc,
		categorySvc: categorySvc,
		bundleSvc:   bundleSvc,
		orderSvc:    orderSvc,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// ─── Stats ────────────────────────────────────────────────────────

func (h *AdminHandler) GetStats(c echo.Context) error {
	totalOrders, _ := h.orderRepo.CountTotal()
	totalRevenue, _ := h.orderRepo.SumRevenue()
	totalProducts, _ := h.productRepo.Count()
	totalCustomers, _ := h.orderRepo.CountCustomers()

	return c.JSON(http.StatusOK, dto.Response{
		Data: dto.StatsResponse{
			TotalOrders:    totalOrders,
			TotalRevenue:   totalRevenue,
			TotalProducts:  totalProducts,
			TotalCustomers: totalCustomers,
		},
	})
}

// ─── Products ─────────────────────────────────────────────────────

func (h *AdminHandler) GetProducts(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	products, total, err := h.productSvc.GetAllAdmin(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data: products, Total: total, Page: page, Limit: limit, TotalPages: totalPages,
	})
}

func (h *AdminHandler) CreateProduct(c echo.Context) error {
	var req dto.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}
	product, err := h.productSvc.Create(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, dto.Response{Data: product, Message: "Product created"})
}

func (h *AdminHandler) UpdateProduct(c echo.Context) error {
	var req dto.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}
	product, err := h.productSvc.Update(c.Param("id"), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Data: product, Message: "Product updated"})
}

func (h *AdminHandler) DeleteProduct(c echo.Context) error {
	if err := h.productSvc.Delete(c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Message: "Product deleted"})
}

// ─── Categories ───────────────────────────────────────────────────

func (h *AdminHandler) GetCategories(c echo.Context) error {
	cats, err := h.categorySvc.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Data: cats})
}

func (h *AdminHandler) CreateCategory(c echo.Context) error {
	var req dto.CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}
	cat, err := h.categorySvc.Create(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, dto.Response{Data: cat, Message: "Category created"})
}

func (h *AdminHandler) UpdateCategory(c echo.Context) error {
	var req dto.CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}
	cat, err := h.categorySvc.Update(c.Param("id"), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Data: cat, Message: "Category updated"})
}

func (h *AdminHandler) DeleteCategory(c echo.Context) error {
	if err := h.categorySvc.Delete(c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Message: "Category deleted"})
}

// ─── Orders ───────────────────────────────────────────────────────

func (h *AdminHandler) GetOrders(c echo.Context) error {
	status := c.QueryParam("status")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	orders, total, err := h.orderSvc.GetAll(status, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data: orders, Total: total, Page: page, Limit: limit, TotalPages: totalPages,
	})
}

func (h *AdminHandler) UpdateOrderStatus(c echo.Context) error {
	var req dto.UpdateOrderStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}
	order, err := h.orderSvc.UpdateStatus(c.Param("id"), req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Data: order, Message: "Status updated"})
}

// ─── Bundles ──────────────────────────────────────────────────────

func (h *AdminHandler) GetBundles(c echo.Context) error {
	bundles, err := h.bundleSvc.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Data: bundles})
}

func (h *AdminHandler) CreateBundle(c echo.Context) error {
	var req dto.CreateBundleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}
	bundle, err := h.bundleSvc.Create(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, dto.Response{Data: bundle, Message: "Bundle created"})
}

func (h *AdminHandler) UpdateBundle(c echo.Context) error {
	var req dto.CreateBundleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}
	bundle, err := h.bundleSvc.Update(c.Param("id"), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Data: bundle, Message: "Bundle updated"})
}

func (h *AdminHandler) DeleteBundle(c echo.Context) error {
	if err := h.bundleSvc.Delete(c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.Response{Message: "Bundle deleted"})
}
