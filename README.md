# Greens & Co. — Backend API

REST API server for Greens & Co., built with **Go + Echo v4** and backed by **PostgreSQL** via GORM. Handles authentication, product catalog, order management, and Midtrans payment integration.

---

## Tech Stack

| Layer | Library / Tool |
|---|---|
| Language | Go 1.24 |
| Framework | [Echo v4](https://echo.labstack.com/) |
| ORM | [GORM v2](https://gorm.io/) + `gorm.io/driver/postgres` |
| Auth | [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) |
| Payment | [Midtrans Go SDK](https://github.com/midtrans/midtrans-go) |
| Validation | [go-playground/validator v10](https://github.com/go-playground/validator) |
| Config | [godotenv](https://github.com/joho/godotenv) |

---

## Project Structure

```
backend/
├── cmd/
│   └── main.go               # Entry point — wires up all layers and starts server
└── internal/
    ├── config/               # Env config loader
    ├── database/             # DB connection, AutoMigrate, seed admin user
    ├── dto/                  # Request / response structs
    ├── handlers/             # HTTP handler layer (Echo)
    │   ├── admin.go
    │   ├── auth.go
    │   ├── bundle.go
    │   ├── category.go
    │   ├── order.go
    │   ├── payment.go
    │   └── product.go
    ├── middleware/           # Auth (JWT) + CORS + security middleware
    ├── models/               # GORM models (User, Product, Category, Bundle, Order, …)
    ├── repositories/         # Database access layer
    └── services/             # Business logic layer
```

---

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL 14+ running locally (or connection string to a remote instance)

### 1. Clone & install dependencies

```bash
git clone <repo-url>
cd backend
go mod download
```

### 2. Configure environment

```bash
cp .env.example .env
```

Edit `.env` with your values:

```dotenv
PORT=8080
DATABASE_URL=postgresql://user:password@localhost:5432/greensco?sslmode=disable
JWT_SECRET=your-secret-key-here-min-32-characters!!
JWT_EXPIRES_IN=7d
MIDTRANS_SERVER_KEY=SB-Mid-server-xxxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxx
MIDTRANS_IS_PRODUCTION=false
CORS_ORIGIN=http://localhost:3000
ADMIN_DEFAULT_PASSWORD=Admin@greensco1
```

> **JWT_SECRET** must be at least 32 characters.  
> **ADMIN_DEFAULT_PASSWORD** must be at least 8 characters and contain at least one number.  
> The database schema is created automatically via GORM `AutoMigrate` on startup.  
> A default admin user (`admin@greensco.com`) is seeded on first run using `ADMIN_DEFAULT_PASSWORD`.

### 3. Run

```bash
go run ./cmd/main.go
```

Or build a binary:

```bash
go build -o greens-co-api ./cmd/main.go
./greens-co-api
```

Server starts at `http://localhost:8080`.

### 4. Health check

```
GET /health
→ 200 { "status": "ok" }
```

---

## API Reference

Base URL: `http://localhost:8080/api`

### Authentication

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/auth/register` | — | Register new user |
| `POST` | `/auth/login` | — | Login, returns JWT |
| `GET` | `/auth/me` | JWT | Get current user profile |

> Auth routes are rate-limited to **10 requests / minute per IP**.

### Public — Catalog

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/products` | List all products |
| `GET` | `/products/:slug` | Get product by slug |
| `GET` | `/categories` | List all categories |
| `GET` | `/bundles` | List all bundles |

### Protected — Orders & Payments

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/orders` | JWT | Place a new order |
| `GET` | `/orders` | JWT | Get current user's orders |
| `POST` | `/payments/initiate` | JWT | Initiate Midtrans payment for an order |
| `POST` | `/payments/notification` | — | Midtrans webhook callback |

### Admin

All admin routes require a valid JWT **and** the `admin` role.

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/admin/stats` | Dashboard stats (revenue, orders, products count) |
| `GET` | `/admin/products` | List all products |
| `POST` | `/admin/products` | Create product |
| `PUT` | `/admin/products/:id` | Update product |
| `DELETE` | `/admin/products/:id` | Delete product |
| `GET` | `/admin/categories` | List all categories |
| `POST` | `/admin/categories` | Create category |
| `PUT` | `/admin/categories/:id` | Update category |
| `DELETE` | `/admin/categories/:id` | Delete category |
| `GET` | `/admin/bundles` | List all bundles |
| `POST` | `/admin/bundles` | Create bundle |
| `PUT` | `/admin/bundles/:id` | Update bundle |
| `DELETE` | `/admin/bundles/:id` | Delete bundle |
| `GET` | `/admin/orders` | List all orders |
| `PATCH` | `/admin/orders/:id` | Update order status |

---

## Security

- **Rate limiting** — auth endpoints capped at 10 req/min per IP (in-memory store, 5-min expiry)
- **Security headers** — `X-XSS-Protection`, `X-Content-Type-Options`, `X-Frame-Options: DENY`, HSTS via Echo `SecureMiddleware`
- **Body limit** — requests capped at 2 MB
- **JWT validation** — secret enforced to be ≥ 32 characters at startup
- **Password policy** — minimum 8 characters, must contain at least one digit
- **CORS** — origin restricted to `CORS_ORIGIN` env var

---

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `PORT` | No | `8080` | Server listen port |
| `DATABASE_URL` | **Yes** | — | PostgreSQL connection string |
| `JWT_SECRET` | **Yes** | — | Secret for signing JWTs (min 32 chars) |
| `JWT_EXPIRES_IN` | No | `7d` | JWT expiry duration |
| `MIDTRANS_SERVER_KEY` | **Yes** | — | Midtrans server key |
| `MIDTRANS_CLIENT_KEY` | **Yes** | — | Midtrans client key |
| `MIDTRANS_IS_PRODUCTION` | No | `false` | Use Midtrans production environment |
| `CORS_ORIGIN` | No | `*` | Allowed CORS origin |
| `ADMIN_DEFAULT_PASSWORD` | No | — | Password for seeded admin account |
