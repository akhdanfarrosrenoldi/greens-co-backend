package main

import (
	"net/http"
	"time"

	"greens-co/backend/internal/config"
	"greens-co/backend/internal/database"
	"greens-co/backend/internal/handlers"
	appmiddleware "greens-co/backend/internal/middleware"
	"greens-co/backend/internal/repositories"
	"greens-co/backend/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// CustomValidator wraps go-playground/validator
type CustomValidator struct {
	v *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}

func main() {
	// Load config
	cfg := config.Load()

	// Connect DB + AutoMigrate + Seed
	db := database.Connect(cfg)

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	bundleRepo := repositories.NewBundleRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Services
	authSvc := services.NewAuthService(userRepo, cfg)
	productSvc := services.NewProductService(productRepo)
	categorySvc := services.NewCategoryService(categoryRepo)
	bundleSvc := services.NewBundleService(bundleRepo)
	orderSvc := services.NewOrderService(orderRepo, productRepo)
	paymentSvc := services.NewPaymentService(orderRepo, cfg)

	// Handlers
	authHandler := handlers.NewAuthHandler(authSvc)
	productHandler := handlers.NewProductHandler(productSvc)
	categoryHandler := handlers.NewCategoryHandler(categorySvc)
	bundleHandler := handlers.NewBundleHandler(bundleSvc)
	orderHandler := handlers.NewOrderHandler(orderSvc)
	paymentHandler := handlers.NewPaymentHandler(paymentSvc)
	adminHandler := handlers.NewAdminHandler(productSvc, categorySvc, bundleSvc, orderSvc, orderRepo, productRepo)

	// Middleware
	authMiddleware := appmiddleware.NewAuthMiddleware(cfg)

	// Echo instance
	e := echo.New()
	e.Validator = &CustomValidator{v: validator.New()}

	// Global middleware
	e.Use(appmiddleware.CORS(cfg))
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.BodyLimit("2M"))
	e.Use(echomiddleware.SecureWithConfig(echomiddleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            3600,
	}))

	// Rate limiter for auth routes (10 req/min per IP)
	authRateLimiter := echomiddleware.RateLimiterWithConfig(echomiddleware.RateLimiterConfig{
		Store: echomiddleware.NewRateLimiterMemoryStoreWithConfig(
			echomiddleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Every(time.Minute / 10),
				Burst:     10,
				ExpiresIn: 5 * time.Minute,
			},
		),
	})

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API group
	api := e.Group("/api")

	// ─── Public routes ────────────────────────────────────────────
	api.POST("/auth/login", authHandler.Login, authRateLimiter)
	api.POST("/auth/register", authHandler.Register, authRateLimiter)
	api.GET("/auth/me", authHandler.Me, authMiddleware.RequireAuth)

	api.GET("/products", productHandler.GetAll)
	api.GET("/products/:slug", productHandler.GetBySlug)
	api.GET("/categories", categoryHandler.GetAll)
	api.GET("/bundles", bundleHandler.GetAll)

	// ─── Protected routes (any authenticated user) ────────────────
	api.POST("/orders", orderHandler.Create, authMiddleware.RequireAuth)
	api.GET("/orders", orderHandler.GetMyOrders, authMiddleware.RequireAuth)
	api.POST("/payments/initiate", paymentHandler.Initiate, authMiddleware.RequireAuth)
	api.POST("/payments/notification", paymentHandler.Notification) // Midtrans webhook

	// ─── Admin routes ─────────────────────────────────────────────
	admin := api.Group("/admin", authMiddleware.RequireAuth, authMiddleware.RequireAdmin)

	admin.GET("/stats", adminHandler.GetStats)

	admin.GET("/products", adminHandler.GetProducts)
	admin.POST("/products", adminHandler.CreateProduct)
	admin.PUT("/products/:id", adminHandler.UpdateProduct)
	admin.DELETE("/products/:id", adminHandler.DeleteProduct)

	admin.GET("/orders", adminHandler.GetOrders)
	admin.PATCH("/orders/:id", adminHandler.UpdateOrderStatus)

	admin.GET("/categories", adminHandler.GetCategories)
	admin.POST("/categories", adminHandler.CreateCategory)
	admin.PUT("/categories/:id", adminHandler.UpdateCategory)
	admin.DELETE("/categories/:id", adminHandler.DeleteCategory)

	admin.GET("/bundles", adminHandler.GetBundles)
	admin.POST("/bundles", adminHandler.CreateBundle)
	admin.PUT("/bundles/:id", adminHandler.UpdateBundle)
	admin.DELETE("/bundles/:id", adminHandler.DeleteBundle)

	// Start server
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
