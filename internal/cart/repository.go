package cart

// Repository is the persistence port for carts, keyed by session id.
type Repository interface {
	Find(id string) (Cart, error)
	Save(c Cart) error
	Delete(id string) error
}
