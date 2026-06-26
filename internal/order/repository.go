package order

// Repository is the persistence port for orders.
type Repository interface {
	FindAll() ([]Order, error)
	FindByID(id int) (Order, error)
	Create(o Order) (Order, error)
	Update(o Order) (Order, error)
	Delete(id int) error
}
