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
- Routes: kebab-case (/product-detail)
- Go structs: PascalCase

## Project structure
drydown/
├── main.go                  # Entry point
├── go.mod
├── go.sum
│
├── internal/                # App code
│   ├── handler/             # HTTP handlers (products, cart, auth...)
│   │   ├── auth.go
│   │   ├── product.go
│   │   ├── cart.go
│   │   └── order.go
│   │
│   ├── model/               # DB queries + structs
│   │   ├── user.go
│   │   ├── product.go
│   │   └── order.go
│   │
│   └── middleware/          # Auth, logging, etc
│       └── auth.go
│
├── templates/               # HTML templates
│   ├── layout.html          # Base layout (nav, footer)
│   ├── index.html
│   ├── product.html
│   ├── cart.html
│   └── partials/            # HTMX fragments
│       ├── product-card.html
│       ├── cart-drawer.html
│       └── search-results.html
│
└── static/                  # CSS, JS, images
    ├── css/
    │   └── app.css          # Tailwind
    ├── js/
    │   ├── htmx.min.js
    │   └── alpine.min.js
    └── img/