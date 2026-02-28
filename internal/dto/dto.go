package dto

// ─── Responses ───────────────────────────────────────────────────

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"totalPages"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Auth response (no data wrapper — token at root)
type AuthResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type UserDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// ─── Auth Requests ────────────────────────────────────────────────

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,containsany=0123456789"`
}

// ─── Product Requests ─────────────────────────────────────────────

type CreateProductRequest struct {
	Name          string                 `json:"name" validate:"required"`
	Slug          string                 `json:"slug"`
	Description   string                 `json:"description"`
	BasePrice     int64                  `json:"basePrice" validate:"required,min=0"`
	Stock         int                    `json:"stock"`
	Image         string                 `json:"image"`
	CategoryID    string                 `json:"categoryId" validate:"required"`
	IsAvailable   bool                   `json:"isAvailable"`
	Badge         *string                `json:"badge"`
	OriginalPrice *int64                 `json:"originalPrice"`
	Variants      []CreateVariantRequest `json:"variants"`
}

type CreateVariantRequest struct {
	ID              *string `json:"id"` // present = update, nil = create new
	Name            string  `json:"name" validate:"required"`
	AdditionalPrice int64   `json:"additionalPrice"`
}

// ─── Category Requests ────────────────────────────────────────────

type CreateCategoryRequest struct {
	Name  string  `json:"name" validate:"required"`
	Slug  string  `json:"slug"`
	Image *string `json:"image"`
}

// ─── Bundle Requests ──────────────────────────────────────────────

type CreateBundleRequest struct {
	Name          string              `json:"name" validate:"required"`
	Slug          string              `json:"slug"`
	Description   string              `json:"description"`
	Price         int64               `json:"price" validate:"required,min=0"`
	OriginalPrice int64               `json:"originalPrice"`
	Image         string              `json:"image"`
	IsPopular     bool                `json:"isPopular"`
	Items         []BundleItemRequest `json:"items"`
}

type BundleItemRequest struct {
	ProductID string `json:"productId" validate:"required"`
	Qty       int    `json:"qty" validate:"required,min=1"`
}

// ─── Order Requests ───────────────────────────────────────────────

type CreateOrderRequest struct {
	Name       string             `json:"name" validate:"required"`
	Phone      string             `json:"phone" validate:"required"`
	Type       string             `json:"type" validate:"required,oneof=DELIVERY PICKUP"`
	Address    *string            `json:"address"`
	Notes      *string            `json:"notes"`
	PickupTime *string            `json:"pickupTime"`
	Items      []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type OrderItemRequest struct {
	ProductID string  `json:"productId" validate:"required"`
	VariantID *string `json:"variantId"`
	Qty       int     `json:"qty" validate:"required,min=1"`
	Price     int64   `json:"price" validate:"required"`
	Notes     *string `json:"notes"`
}

// ─── Admin Requests ───────────────────────────────────────────────

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=PENDING PAID PROCESSING READY ON_DELIVERY COMPLETED CANCELLED"`
}

// ─── Stats Response ───────────────────────────────────────────────

type StatsResponse struct {
	TotalOrders    int64 `json:"totalOrders"`
	TotalRevenue   int64 `json:"totalRevenue"`
	TotalProducts  int64 `json:"totalProducts"`
	TotalCustomers int64 `json:"totalCustomers"`
}

// ─── Payment ──────────────────────────────────────────────────────

type InitiatePaymentRequest struct {
	OrderID string `json:"orderId" validate:"required"`
}

type PaymentResponse struct {
	PaymentURL string `json:"paymentUrl"`
}
