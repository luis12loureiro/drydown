package product

// Repository is the persistence port for products. It is implemented by an
// adapter (memory.go now, sqlite.go later) and consumed by the Service.
type Repository interface {
	FindAll() ([]Product, error)
	FindByFamily(fam string) ([]Product, error)
	FindByID(id int) (Product, error)
	Create(p Product) (Product, error)
	Update(p Product) (Product, error)
	Delete(id int) error
}
