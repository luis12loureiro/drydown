package order

// Repository is the persistence port for orders.
type Repository interface {
	FindAll() ([]Order, error)
	FindByUUID(uuid string) (Order, error)
	Create(o Order) (Order, error)
	Update(o Order) (Order, error)
	Delete(uuid string) error
}
