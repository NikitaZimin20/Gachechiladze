package product

import "time"

type Product struct {
	ID          uint32
	Title       string  `schema:"title"`
	Price       float32 `schema:"price"`
	Tag         string  `schema:"tag"`
	Type        string  `schema:"type"`
	Description string  `schema:"description"`
	Count       uint32
	CreatedAt   time.Time `schema:"created_at"`
	Views       uint32
	ImageURL    string
}

type OrderProduct struct {
	UserID      uint32
	ID          uint32
	Title       string
	Description string
	Price       float32
	ImageURL    string
}

type Order struct {
	Products   []*OrderProduct
	TotalPrice float32
}

type DeliveryOrder struct {
	UserID        uint32
	ProductName   string
	Price         float32
	DepartureDate time.Time
	LandingDate   time.Time
	Status        string
}

type ProductRepo interface {
	GetAll(orderBy string) ([]*Product, error)
	GetByID(id uint32) (*Product, error)
	Add(product *Product) (uint32, error)
	Update(newProduct *Product) (bool, error)
	Delete(id uint32) (bool, error)

	GetOrdersByID(id uint32) (*Order, error)
	AddOrder(id uint32, product *Product) (uint32, error)
	DeleteOrder(id, prodID uint32) (bool, error)

	RegisterOrder(userID uint32, order []*OrderProduct) (uint32, error)
	GetDeliveryOrdersByID(userID uint32) ([]*DeliveryOrder, error)
	GetRelated(typ string, limit int) ([]*Product, error)
}
