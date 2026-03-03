# E‑Commerce API Service (Intermediate)

A layered REST API that implements a minimal e‑commerce backend with authentication, product catalog, cart, checkout, and orders. It uses SQLite for persistence, JWT for auth, chi for routing, and Stripe for payment intents.

This project exposes:
- A JSON API under /api/v1
- A lightweight server‑rendered UI that now talks to the JSON API exclusively via fetch

## Highlights
- Clean layering: handlers → services → repositories → database
- JWT auth with role‑based access (admin/customer)
- SQLite with pragmatic migrations and WAL mode
- Atomic checkout with a single DB transaction
- Structured middleware: logging, auth, admin‑only
- Stripe PaymentIntent integration (test mode)
- OpenAPI (Swagger) docs and in‑app Swagger UI

## Architecture
- config: reads environment variables into a typed Config
- database: opens SQLite with WAL and runs idempotent migrations; seeds admin
- models: request/response and domain types plus reusable errors
- repository: SQL queries; no business logic
- services: business rules and transactions; composes repositories
- handlers: HTTP I/O, JSON parsing, error→status mapping (API); plus page rendering for the HTML skeletons
- router: all routes; applies middleware and groups
- middleware: request logging, JWT auth, admin gate

## Tech Stack
- Go, chi/v5, golang‑jwt/jwt, bcrypt
- SQLite (github.com/mattn/go‑sqlite3)
- Stripe (stripe‑go v84)
- Swagger tooling (swag, http‑swagger)

## Directory Layout
```
config/        env loading
database/      sqlite open, migrations, seed
handlers/      HTTP handlers (auth, product, cart, order, pages)
middleware/    logging, auth, admin
models/        DTOs and domain models
repository/    SQL data access
router/        chi routes & wiring
services/      business logic (auth, product, cart, order, payment)
templates/     thin HTML pages; JS fetches call JSON API
docs/          generated OpenAPI documentation (swagger.json/.yaml)
main.go        composition root
```

## Requirements
- Go 1.21+ (tested against newer Go)
- CGO enabled for mattn/go‑sqlite3
  - macOS: install Xcode Command Line Tools (`xcode-select --install`)
  - Linux: install gcc/clang and sqlite dev headers
- A Stripe test secret key for checkout (optional if you don’t call /checkout)

## Configuration
Environment variables with defaults:
- JWT_SECRET: default "change-me-in-production"
- STRIPE_SECRET_KEY: default "" (checkout fails if unset)
- DB_PATH: default "ecommerce.db" (stored under ./data/)
- PORT: default "8080"
- ADMIN_EMAIL: default "admin@jaygaha.com.np"
- ADMIN_PASSWORD: default "admin123"

## Setup
1) Ensure CGO is enabled (required by go‑sqlite3):
```
export CGO_ENABLED=1
```
2) Create the data directory (SQLite file lives here):
```
mkdir -p data
```
3) Install dependencies:
```
go mod download
```
4) (Optional) Set your env:
```
export JWT_SECRET=supersecret
export STRIPE_SECRET_KEY=sk_test_...
export PORT=8080
```
5) Run:
```
go run main.go
```
On boot the app applies migrations and seeds an admin user (ADMIN_EMAIL / ADMIN_PASSWORD).

## API Documentation (Swagger/OpenAPI)
- Live Swagger UI: http://localhost:8080/swagger
- BasePath: /api/v1 (applies to every API operation)
- Regenerate docs after annotation changes:
```
go run github.com/swaggo/swag/cmd/swag@latest init -g main.go -o docs
```
This writes docs/docs.go, docs/swagger.json, and docs/swagger.yaml.

## API
Base URL: /api/v1

Public
- POST /auth/register
- POST /auth/login
- GET  /products
- GET  /products/{id}

Authenticated (Bearer token)
- GET    /cart
- POST   /cart/items
- PUT    /cart/items/{productId}
- DELETE /cart/items/{productId}
- POST   /checkout
- GET    /orders
- GET    /orders/{id}

Admin (Bearer token with role=admin)
- POST   /admin/products
- PUT    /admin/products/{id}
- DELETE /admin/products/{id}

## Frontend (API‑Only)
The server renders simple HTML pages from templates/, but all data‑bearing interactions now use the JSON API via fetch with a Bearer token:
- Sign In / Sign Up: Pages call /api/v1/auth/login and /api/v1/auth/register and store the JWT in localStorage.
- Navigation: The base layout decodes the JWT payload to toggle Admin and Sign Out visibility, and shows Sign In when no token is present.
- Products: Home and product pages fetch product lists/details via /api/v1 endpoints.
- Cart: Pages call /api/v1/cart and /api/v1/cart/items (POST/PUT/DELETE) and render results dynamically.
- Checkout: Triggers POST /api/v1/checkout and redirects to Orders.

Notes:
- Because the token is stored in localStorage, all API requests are sent with Authorization: Bearer <token>.
- If you prefer HttpOnly cookies for the session, you can add a small “cookie bridge” or switch middleware to also read JWTs from cookies.

## Example Workflow (curl)
Register and login:
```
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"you@example.com","password":"password123"}'

TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"you@example.com","password":"password123"}' | jq -r .token)
```

Create a product (admin only; use seeded admin or promote manually):
```
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"'"$ADMIN_EMAIL"'","password":"'"$ADMIN_PASSWORD"'"}' | jq -r .token)

curl -X POST http://localhost:8080/api/v1/admin/products \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"T‑Shirt","description":"Soft cotton","price":1999,"stock":10,"image_url":"https://example.com/t.png"}'
```

Browse and add to cart:
```
curl http://localhost:8080/api/v1/products
curl -X POST http://localhost:8080/api/v1/cart/items \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"product_id":1,"quantity":2}'
curl http://localhost:8080/api/v1/cart -H "Authorization: Bearer $TOKEN"
```

Checkout (requires STRIPE_SECRET_KEY):
```
curl -X POST http://localhost:8080/api/v1/checkout \
  -H "Authorization: Bearer $TOKEN"
```
Response includes client_secret and stripe_payment_id.

## Design Notes
- Handlers are transport‑level only; services host business rules
- Repositories parameterize all queries; no string concatenation with inputs
- Checkout uses a single transaction: create order → decrement stock → order items → clear cart → commit
- WAL and foreign_keys pragmas improve read concurrency and integrity
- The HTML UI is a thin layer over the API; it can be progressively enhanced or replaced with a SPA without changing backend contracts

## Troubleshooting
- sql: unknown driver "sqlite3": ensure you imported and built the driver and have CGO enabled:
  - Dependency in go.mod: github.com/mattn/go-sqlite3
  - Environment: `export CGO_ENABLED=1`
  - macOS: install Xcode CLT
- cannot open data/...: ensure `mkdir -p data`
- 401 unauthorized: send `Authorization: Bearer <token>` header
- 403 forbidden on admin endpoints: login as admin; defaults set by ADMIN_EMAIL / ADMIN_PASSWORD
- Checkout errors about Stripe: set STRIPE_SECRET_KEY to a Stripe test key
- Swagger UI not loading: confirm /swagger is routed and docs/ exists; regenerate with `swag init`
- Admin menu hidden: ensure you logged in via /auth so a JWT is present in localStorage and role=admin in the token payload

## Project Link

- [E-Commerce API](https://roadmap.sh/projects/ecommerce-api)

## Acknowledgments

- Part of the Go programming language learning roadmap projects
- Created by [jaygaha](https://github.com/jaygaha)
