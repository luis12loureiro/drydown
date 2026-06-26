// Package user holds the customer domain. Authentication is intentionally out
// of scope for the first deploy; this is plain customer data only.
package user

import (
	"errors"
	"time"
)

// ErrNotFound is returned when a user cannot be located by id.
var ErrNotFound = errors.New("user: not found")

// ErrInvalid is returned when a user fails validation.
var ErrInvalid = errors.New("user: invalid")

// User is a customer record. Kept minimal for the first deploy.
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Valid reports whether the user carries the minimum required fields.
func (u User) Valid() bool {
	return u.Name != "" && u.Email != ""
}
