// Package order holds the checkout domain: an Order and its line items.
package order

import (
	"errors"
	"time"
)

// ErrNotFound is returned when an order cannot be located by its uuid.
var ErrNotFound = errors.New("order: not found")

// ErrEmpty is returned when trying to place an order with no items.
var ErrEmpty = errors.New("order: no items")

// Order status values.
const (
	StatusPending   = "pending"
	StatusPaid      = "paid"
	StatusShipped   = "shipped"
	StatusCancelled = "cancelled"
)

// Item is a single purchased line, captured at the time of checkout.
type Item struct {
	ProductUUID string
	Name        string
	ML          int
	Price       int
	Qty         int
}

// Order is a placed order. Attributes are kept minimal for the first deploy.
type Order struct {
	UUID      string
	Items     []Item
	Total     int
	Status    string
	CreatedAt time.Time
}
