package product

import "time"

type Product struct {
	ID            uint32
	Title         string  `schema:"title"`
	Price         float32 `schema:"price"`
	Tag           string  `schema:"tag"`
	Type          string  `schema:"type"`
	Description   string  `schema:"description"`
	FactoryCount  uint32
	CreatedAt     time.Time `schema:"created_at"`
	PurchaseCount uint32
	Views         uint32
	ImageURL      string
}

type Basket struct {
	UserID     uint32
	TotalPrice float32
	TotalCount uint32
	Products   []*Product
}

type LandingInfo struct {
	ProductName   string
	Price         float32
	DepartureDate time.Time
	LandingDate   time.Time
	Status        string
}

type History struct {
	UserID   uint32
	Landings []*LandingInfo
}

type ProductRepo interface {
	GetBasketByID(id uint32) (*Basket, error)
	AddBasket(id uint32) (uint32, error)
	AddProductToBasket(id uint32, product *Product) (uint32, error)
	DeleteProductFromBasket(id uint32, prodID uint32) (bool, error)
	GetAll(orderBy string) ([]*Product, error)
	GetByID(id uint32) (*Product, error)
	Add(product *Product) (uint32, error)
	Update(newProduct *Product) (bool, error)
	Delete(id uint32) (bool, error)
	RegisterOrder(userID uint32, products []*Product) (uint32, error)
	GetOrders(userID uint32) ([]*LandingInfo, error)
	GetRelated(typ string, limit int) ([]*Product, error)
}
