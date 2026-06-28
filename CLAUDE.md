# Drydown

## What this is
An e-commerce store for fragrance decants, built with Go + HTMX + Tailwind.
Later will include our own fragrance line.

## Stack
- Backend: Go, Chi router
- Frontend: HTMX, Alpine.js, Tailwind CSS
- Database: SQLite and later PostgreSQL
- Payments: Stripe

## Project rules
- No JS frameworks, HTMX only for interactivity
- Handlers must be thin — no business logic, only call models + render templates
- No raw SQL in handlers, always go through service and then model layer
- All HTMX partials live in templates/partials/
- Use httpOnly cookies for auth with sessions set on the DB
- Name of the go package must match the name of the directory that contains the package
- Go package names should describe what the package  (e.g dont call a package util, call it format)

## Naming conventions
- Files: snake_case
- Go files: smallcase
- Routes: kebab-case (/product-detail)
- Go structs: PascalCase

## Project structure
drydown/
├── main.go
├── go.mod
├── go.sum
│
├── internal/
│   │
│   ├── product/
│   │   ├── product.go           # Domain struct + rules
│   │   ├── repository.go        # Interface (port)
│   │   ├── service.go           # Business logic
│   │   ├── handler.go           # HTTP handler
│   │   └── sqlite.go            # DB implementation (adapter)
│   │
│   ├── order/
│   │   ├── order.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── sqlite.go
│   │
│   ├── cart/
│   │   ├── cart.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── sqlite.go
│   │
│   ├── user/
│   │   ├── user.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── sqlite.go
│   │
│   └── payment/
│       ├── payment.go           # Interface
│       └── stripe.go            # Stripe implementation
│
├── templates/
│   ├── layout.html
│   ├── index.html
│   └── partials/
│       ├── product-card.html
│       ├── cart-drawer.html
│       └── search-results.html
│
└── static/
    ├── css/
    ├── js/
    └── img/