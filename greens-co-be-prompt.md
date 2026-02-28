# Claude Code Prompt — Greens & Co. Backend (Go + Echo)


```
Build a complete REST API backend for "Greens & Co." — a healthy F&B e-commerce platform.

The frontend is already built in Next.js. This backend must implement exactly 28 endpoints
as specified. Do NOT deviate from the response format — the frontend depends on it.

---

## TECH STACK
- Language  : Go 1.22+
- Framework : Echo v4 (github.com/labstack/echo/v4)
- ORM       : GORM v2 (gorm.io/gorm) + PostgreSQL driver (gorm.io/driver/postgres)
- Auth      : JWT (github.com/golang-jwt/jwt/v5)
- Password  : bcrypt (golang.org/x/crypto/bcrypt)
- Config    : godotenv (github.com/joho/godotenv)
- Payment   : Midtrans Go SDK (github.com/midtrans/midtrans-go)
- CORS      : Echo built-in middleware
- Validator : go-playground/validator/v10

---

## FOLDER STRUCTURE

```
backend/
├── cmd/
│   └── main.go                  ← entry point
├── internal/
│   ├── config/
│   │   └── config.go            ← load .env, return Config struct
│   ├── database/
│   │   └── database.go          ← GORM connection + AutoMigrate
│   ├── middleware/
│   │   ├── auth.go              ← JWT middleware (RequireAuth, RequireAdmin)
│   │   └── cors.go              ← CORS config
│   ├── models/
│   │   └── models.go            ← all GORM models in one file
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── product.go
│   │   ├── category.go
│   │   ├── bundle.go
│   │   ├── order.go
│   │   ├── payment.go
│   │   └── admin.go
│   ├── services/
│   │   ├── auth.go
│   │   ├── product.go
│   │   ├── category.go
│   │   ├── bundle.go
│   │   ├── order.go
│   │   └── payment.go
│   ├── repositories/
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── category.go
│   │   ├── bundle.go
│   │   └── order.go
│   └── dto/
│       └── dto.go               ← all request/response DTOs
├── .env
├── .env.example
├── go.mod
└── go.sum
```

---

## MODELS (`internal/models/models.go`)

```go
type User struct {
    ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    Role      string    `gorm:"default:'CUSTOMER'"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Category struct {
    ID       string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name     string     `gorm:"not null"`
    Slug     string     `gorm:"uniqueIndex;not null"`
    Image    *string
    Products []Product  `gorm:"foreignKey:CategoryID"`
}

type Product struct {
    ID            string           `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name          string           `gorm:"not null"`
    Slug          string           `gorm:"uniqueIndex;not null"`
    Description   string
    BasePrice     int64            `gorm:"not null"`
    OriginalPrice *int64
    Image         string
    Stock         int              `gorm:"default:0"`
    IsAvailable   bool             `gorm:"default:true"`
    Badge         *string          // "bestseller" | "new" | "promo" | null
    Rating        *float64
    ReviewCount   *int
    CategoryID    string           `gorm:"type:uuid;not null"`
    Category      Category         `gorm:"foreignKey:CategoryID"`
    Variants      []ProductVariant `gorm:"foreignKey:ProductID"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type ProductVariant struct {
    ID              string  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    ProductID       string  `gorm:"type:uuid;not null"`
    Name            string  `gorm:"not null"`
    AdditionalPrice int64   `gorm:"default:0"`
}

type Bundle struct {
    ID            string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name          string       `gorm:"not null"`
    Slug          string       `gorm:"uniqueIndex;not null"`
    Description   string
    Price         int64        `gorm:"not null"`
    OriginalPrice int64        `gorm:"not null"`
    Image         string
    IsPopular     bool         `gorm:"default:false"`
    Items         []BundleItem `gorm:"foreignKey:BundleID"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type BundleItem struct {
    ID        string  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    BundleID  string  `gorm:"type:uuid;not null"`
    ProductID string  `gorm:"type:uuid;not null"`
    Product   Product `gorm:"foreignKey:ProductID"`
    Qty       int     `gorm:"not null"`
}

type Order struct {
    ID            string      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    UserID        string      `gorm:"type:uuid;not null"`
    User          User        `gorm:"foreignKey:UserID"`
    Status        string      `gorm:"default:'PENDING'"`
    Type          string      `gorm:"not null"` // DELIVERY | PICKUP
    TotalPrice    int64       `gorm:"not null"`
    Name          string
    Phone         string
    Address       *string
    PickupTime    *string
    Notes         *string
    PaymentStatus string      `gorm:"default:'UNPAID'"`
    MidtransID    *string
    Items         []OrderItem `gorm:"foreignKey:OrderID"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type OrderItem struct {
    ID        string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    OrderID   string          `gorm:"type:uuid;not null"`
    ProductID string          `gorm:"type:uuid;not null"`
    Product   Product         `gorm:"foreignKey:ProductID"`
    VariantID *string         `gorm:"type:uuid"`
    Variant   *ProductVariant `gorm:"foreignKey:VariantID"`
    Qty       int             `gorm:"not null"`
    Price     int64           `gorm:"not null"`
    Notes     *string
}
```

---

## ALL 28 ENDPOINTS

### Routes setup in `cmd/main.go`

```go
e := echo.New()

// Middleware
e.Use(middleware.CORS())
e.Use(middleware.Logger())
e.Use(middleware.Recover())

api := e.Group("/api")

// Public
api.POST("/auth/login", authHandler.Login)
api.POST("/auth/register", authHandler.Register)
api.GET("/auth/me", authHandler.Me, authMiddleware.RequireAuth)

api.GET("/products", productHandler.GetAll)
api.GET("/products/:slug", productHandler.GetBySlug)
api.GET("/categories", categoryHandler.GetAll)
api.GET("/bundles", bundleHandler.GetAll)

// Protected (any logged-in user)
api.POST("/orders", orderHandler.Create, authMiddleware.RequireAuth)
api.GET("/orders", orderHandler.GetMyOrders, authMiddleware.RequireAuth)
api.POST("/payments/initiate", paymentHandler.Initiate, authMiddleware.RequireAuth)
api.POST("/payments/notification", paymentHandler.Notification) // Midtrans webhook, no auth

// Admin only
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
```

---

## RESPONSE FORMAT CONVENTIONS

**CRITICAL — frontend depends on these exact shapes:**

```go
// Single resource
type Response struct {
    Data    interface{} `json:"data"`
    Message string      `json:"message,omitempty"`
}

// Paginated list
type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    Limit      int         `json:"limit"`
    TotalPages int         `json:"totalPages"`
}

// Error
type ErrorResponse struct {
    Message string      `json:"message"`
    Errors  interface{} `json:"errors,omitempty"`
}

// Auth (no wrapper — token at root level)
type AuthResponse struct {
    Token string   `json:"token"`
    User  UserDTO  `json:"user"`
}
```

---

## KEY ENDPOINT BEHAVIORS

### GET /products (with filters)
Accept query params: `cat`, `sort`, `min`, `max`, `available`, `diet`, `search`, `page`, `limit`
- `cat` → filter by category slug
- `sort` → popular (reviewCount DESC) | newest (createdAt DESC) | price_asc | price_desc
- `min`/`max` → filter by basePrice range
- `available` → filter by isAvailable (default true)
- `search` → ILIKE %name%
- `page` default 1, `limit` default 9
- Preload Category and Variants

### POST /auth/register
- Hash password with bcrypt (cost 12)
- Return AuthResponse (token + user), status 201

### POST /auth/login
- Verify bcrypt password
- Return AuthResponse with JWT (expires JWT_EXPIRES_IN from env)

### JWT Claims
```go
type JWTClaims struct {
    UserID string `json:"userId"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}
```

### RequireAuth middleware
- Extract Bearer token from Authorization header
- Validate JWT, set `userID` and `role` in echo context

### RequireAdmin middleware
- Check `role` from echo context == "ADMIN"
- Return 403 if not

### POST /orders
- Calculate totalPrice from items (price × qty) + delivery fee (10000 if DELIVERY, 0 if PICKUP)
- Decrement product stock for each item
- Return created Order with items preloaded

### POST /payments/initiate
- Create Midtrans Snap transaction
- Return { data: { paymentUrl: string } }
- Use sandbox if MIDTRANS_IS_PRODUCTION=false

### POST /payments/notification (Midtrans webhook)
- Verify notification signature
- Update Order status: if transaction_status=settlement → status=PAID, paymentStatus=PAID

### GET /admin/stats
```go
// Count total orders, sum revenue (status != CANCELLED), count products, count unique users with orders
```

### Slug auto-generation
For POST /admin/products and POST /admin/bundles, if slug is empty:
```go
slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
// remove special chars, ensure uniqueness by appending -2, -3 if exists
```

---

## DTO EXAMPLES (`internal/dto/dto.go`)

```go
// Auth
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// Product
type CreateProductRequest struct {
    Name          string                   `json:"name" validate:"required"`
    Slug          string                   `json:"slug"`
    Description   string                   `json:"description"`
    BasePrice     int64                    `json:"basePrice" validate:"required,min=0"`
    Stock         int                      `json:"stock"`
    Image         string                   `json:"image"`
    CategoryID    string                   `json:"categoryId" validate:"required"`
    IsAvailable   bool                     `json:"isAvailable"`
    Badge         *string                  `json:"badge"`
    OriginalPrice *int64                   `json:"originalPrice"`
    Variants      []CreateVariantRequest   `json:"variants"`
}

type CreateVariantRequest struct {
    ID              *string `json:"id"` // present = update, nil = create new
    Name            string  `json:"name" validate:"required"`
    AdditionalPrice int64   `json:"additionalPrice"`
}

// Order
type CreateOrderRequest struct {
    Name       string              `json:"name" validate:"required"`
    Phone      string              `json:"phone" validate:"required"`
    Type       string              `json:"type" validate:"required,oneof=DELIVERY PICKUP"`
    Address    *string             `json:"address"`
    Notes      *string             `json:"notes"`
    PickupTime *string             `json:"pickupTime"`
    Items      []OrderItemRequest  `json:"items" validate:"required,min=1"`
}

type OrderItemRequest struct {
    ProductID string  `json:"productId" validate:"required"`
    VariantID *string `json:"variantId"`
    Qty       int     `json:"qty" validate:"required,min=1"`
    Price     int64   `json:"price" validate:"required"`
    Notes     *string `json:"notes"`
}

// Admin
type UpdateOrderStatusRequest struct {
    Status string `json:"status" validate:"required,oneof=PENDING PAID PROCESSING READY ON_DELIVERY COMPLETED CANCELLED"`
}
```

---

## SEED DATA

After AutoMigrate, seed this data if tables are empty:

```go
// Categories
categories := []Category{
    {Name: "Salad",     Slug: "salad"},
    {Name: "Rice Bowl", Slug: "rice-bowl"},
    {Name: "Drinks",    Slug: "drinks"},
    {Name: "Snack",     Slug: "snack"},
}

// Products (use category IDs from above)
products := []Product{
    {Name: "Garden Fresh Salad",   Slug: "garden-fresh-salad",   BasePrice: 35000, Badge: ptr("bestseller"), Rating: ptr(4.9), ReviewCount: ptr(128), Stock: 10, IsAvailable: true, Image: "https://images.unsplash.com/photo-1512621776951-a57141f2eefd?w=400&q=80", Description: "Mixed greens, cherry tomato, cucumber, house dressing"},
    {Name: "Teriyaki Chicken Bowl", Slug: "teriyaki-chicken-bowl", BasePrice: 45000, Badge: ptr("new"),        Rating: ptr(4.8), ReviewCount: ptr(94),  Stock: 8,  IsAvailable: true, Image: "https://images.unsplash.com/photo-1547592180-85f173990554?w=400&q=80", Description: "Steamed rice, grilled chicken, teriyaki sauce, sesame"},
    {Name: "Green Detox Juice",     Slug: "green-detox-juice",     BasePrice: 28000, Badge: nil,               Rating: ptr(4.7), ReviewCount: ptr(76),  Stock: 15, IsAvailable: true, Image: "https://images.unsplash.com/photo-1610970881699-44a5587cabec?w=400&q=80", Description: "Spinach, apple, ginger, lemon, cucumber blend"},
    {Name: "Overnight Oats",        Slug: "overnight-oats",        BasePrice: 32000, Badge: nil,               Rating: ptr(4.8), ReviewCount: ptr(53),  Stock: 12, IsAvailable: true, Image: "https://images.unsplash.com/photo-1563805042-7684c019e1cb?w=400&q=80", Description: "Rolled oats, chia seeds, almond milk, mixed berries"},
    {Name: "Quinoa Power Bowl",     Slug: "quinoa-power-bowl",     BasePrice: 48000, Badge: ptr("promo"),      Rating: ptr(4.9), ReviewCount: ptr(112), Stock: 6,  IsAvailable: true, Image: "https://images.unsplash.com/photo-1540420773420-3366772f4999?w=400&q=80", Description: "Quinoa, roasted veggies, tahini dressing, seeds", OriginalPrice: ptr(int64(55000))},
    {Name: "Açaí Bowl",             Slug: "acai-bowl",             BasePrice: 52000, Badge: nil,               Rating: ptr(4.9), ReviewCount: ptr(145), Stock: 0,  IsAvailable: false, Image: "https://images.unsplash.com/photo-1490645935967-10de6ba17061?w=400&q=80", Description: "Blended açaí, granola, fresh fruits, honey drizzle"},
}

// Admin user
admin := User{
    Name:     "Admin Greens",
    Email:    "admin@greensco.id",
    Password: hashPassword("admin123"), // bcrypt
    Role:     "ADMIN",
}

// Bundles
bundles := []Bundle{
    {Name: "Healthy Starter", Slug: "healthy-starter", Price: 79000,  OriginalPrice: 95000,  IsPopular: false, Image: "https://images.unsplash.com/photo-1540420773420-3366772f4999?w=600&q=80", Description: "Perfect for a light & nutritious meal."},
    {Name: "Full Day Pack",   Slug: "full-day-pack",   Price: 125000, OriginalPrice: 140000, IsPopular: true,  Image: "https://images.unsplash.com/photo-1498837167922-ddd27525d352?w=600&q=80", Description: "Complete nutrition for your entire day."},
    {Name: "Family Pack",     Slug: "family-pack",     Price: 215000, OriginalPrice: 256000, IsPopular: false, Image: "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?w=600&q=80", Description: "Feed the whole family with goodness."},
}
```

---

## ENV FILE (`.env`)

```
PORT=8080
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/greensco?sslmode=disable
JWT_SECRET=greensco-secret-key-change-in-production
JWT_EXPIRES_IN=7d
MIDTRANS_SERVER_KEY=SB-Mid-server-xxxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxx
MIDTRANS_IS_PRODUCTION=false
CORS_ORIGIN=http://localhost:3000
```

---

## INIT COMMANDS

```bash
# Di dalam folder backend/
go mod init greens-co/backend
go get github.com/labstack/echo/v4
go get github.com/labstack/echo/v4/middleware
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/joho/godotenv
go get github.com/go-playground/validator/v10
go get github.com/midtrans/midtrans-go

# Jalankan
go run cmd/main.go
```

---

## IMPORTANT NOTES

1. Use Repository Pattern strictly: Handler → Service → Repository → DB
2. Never return password field in any response
3. All UUIDs use PostgreSQL gen_random_uuid() — enable with: `CREATE EXTENSION IF NOT EXISTS pgcrypto;`
4. CORS must allow: origin=http://localhost:3000, methods=GET/POST/PUT/PATCH/DELETE, headers=Content-Type/Authorization
5. Validate all request bodies using go-playground/validator
6. Return proper HTTP status codes (201 for created, 400 for validation error, 401 for unauthorized, 403 for forbidden, 404 for not found)
7. Preload associations where needed (Product → Category + Variants, Order → Items → Product + Variant)
8. Use `c.Bind()` + `c.Validate()` pattern in Echo handlers
9. For Midtrans notification webhook — verify using MIDTRANS_SERVER_KEY signature
10. Product stock decrement must be in a DB transaction with order creation
```

---

## CARA PAKAI

```bash
# Di folder greens-co/
mkdir backend && cd backend

# Buka Claude Code
claude

# Paste prompt di atas
```

Setelah generate:
```bash
# Setup PostgreSQL database dulu
createdb greensco

# Copy .env.example ke .env dan isi nilai yang benar
cp .env.example .env

# Jalankan
go run cmd/main.go
```

Server akan jalan di `http://localhost:8080`

---

## TEST KONEKSI FE + BE

Setelah BE jalan, buka FE di `http://localhost:3000`:
1. Cek landing page — produk harus muncul dari DB
2. Register akun baru
3. Tambah produk ke cart → checkout
4. Login sebagai `admin@greensco.id` / `admin123` → cek `/admin`
