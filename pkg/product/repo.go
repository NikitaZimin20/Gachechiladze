package product

import (
	"errors"
	"sort"
	"time"
)

type ProductMemoryRepository struct {
	data           []*Product
	orders         []*OrderProduct
	deliveryOrders []*DeliveryOrder
	lastID         uint32
}

func NewMemoryRepo() *ProductMemoryRepository {
	return &ProductMemoryRepository{
		data: []*Product{
			{
				ID:          1,
				Title:       "Хайлайтер для губ",
				Price:       42.33,
				Description: "Описание1",
				Count:       5,
				Type:        "косметика",
				Tag:         "Хит",
				ImageURL:    "/static/images/Highlighter.png",
			},
			{
				ID:          2,
				Title:       "Подводка",
				Price:       95.00,
				Description: "Описание2",
				Count:       5,
				Type:        "косметика",
				Tag:         "Sale",
				ImageURL:    "/static/images/fit_black.png",
			},
			{
				ID:          3,
				Title:       "Набор из блесков",
				Price:       156.00,
				Description: "Описание3",
				Count:       5,
				Type:        "косметика",
				Tag:         "Sale",
				ImageURL:    "/static/images/fit_mini.png",
			},
			{
				ID:          4,
				Title:       "Помада Клубничная",
				Price:       180,
				Description: "Описание4",
				Count:       5,
				Type:        "косметика",
				Tag:         "Хит",
				ImageURL:    "/static/images/lip.png",
			},
			{
				ID:          5,
				Title:       "Помада карандаш",
				Price:       200,
				Description: "Описание5",
				Count:       5,
				Type:        "косметика",
				Tag:         "Хит",
				ImageURL:    "/static/images/lip_array.png",
			},
			{
				ID:          6,
				Title:       "Помада Матовая",
				Price:       249,
				Description: "Описание6",
				Count:       5,
				Type:        "косметика",
				Tag:         "Хит",
				ImageURL:    "/static/images/lip_pink.png",
			},
			{
				ID:          7,
				Title:       "Ингаляторы Нирдош",
				Price:       50000000,
				Description: "Описание7",
				Count:       5,
				Type:        "косметика",
				Tag:         "Хит",
				ImageURL:    "/static/images/ing.png",
			},
			{
				ID:          8,
				Title:       "Паста для шугаринга",
				Price:       600,
				Description: "Описание8",
				Count:       5,
				Type:        "косметика",
				Tag:         "Хит",
				ImageURL:    "/static/images/Paste.png",
			},
		},
		orders:         []*OrderProduct{},
		deliveryOrders: []*DeliveryOrder{},
		lastID:         8,
	}
}

var (
	ErrNoBasket = errors.New("no basket found")
)

func (repo *ProductMemoryRepository) GetAll(orderBy string) ([]*Product, error) {
	data := repo.data
	switch orderBy {
	case "views":
		sort.Slice(data, func(i, j int) bool {
			return data[i].Views > data[j].Views
		})
	case "created_at":
		sort.Slice(data, func(i, j int) bool {
			return data[i].CreatedAt.Unix() > data[j].CreatedAt.Unix()
		})
	default:
		break
	}
	return data, nil
}

func (repo *ProductMemoryRepository) GetByID(id uint32) (*Product, error) {
	for _, item := range repo.data {
		if item.ID == id {
			return item, nil
		}
	}
	return &Product{}, nil
}

func (repo *ProductMemoryRepository) Add(product *Product) (uint32, error) {
	repo.lastID++
	product.ID = repo.lastID
	repo.data = append(repo.data, product)
	return repo.lastID, nil
}

func (repo *ProductMemoryRepository) Update(newProduct *Product) (bool, error) {
	for _, item := range repo.data {
		if item.ID != newProduct.ID {
			continue
		}
		item.Views = newProduct.Views
		item.Description = newProduct.Description
		item.Price = newProduct.Price
		return true, nil
	}
	return false, nil
}

func (repo *ProductMemoryRepository) Delete(id uint32) (bool, error) {
	i := -1
	for idx, item := range repo.data {
		if item.ID != id {
			continue
		}
		i = idx
	}
	if i < 0 {
		return false, nil
	}

	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil // or the zero value of T
	repo.data = repo.data[:len(repo.data)-1]

	return true, nil
}

func (repo *ProductMemoryRepository) GetOrdersByID(id uint32) (*Order, error) {
	order := &Order{}
	for _, ord := range repo.orders {
		if ord.UserID == id {
			order.Products = append(order.Products, ord)
			order.TotalPrice += ord.Price
		}
	}
	return order, nil
}

func (repo *ProductMemoryRepository) AddOrder(id uint32, product *Product) (uint32, error) {
	order := &OrderProduct{
		UserID:      id,
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		ImageURL:    product.ImageURL,
	}
	repo.orders = append(repo.orders, order)
	return order.ID, nil
}

func (repo *ProductMemoryRepository) DeleteOrder(id, prodID uint32) (bool, error) {
	i := -1
	for idx, item := range repo.orders {
		if item.ID != prodID || item.UserID != id {
			continue
		}
		i = idx
	}
	if i < 0 {
		return false, nil
	}

	if i < len(repo.orders)-1 {
		copy(repo.orders[i:], repo.orders[i+1:])
	}
	repo.orders[len(repo.orders)-1] = &OrderProduct{} // or the zero value of T
	repo.orders = repo.orders[:len(repo.orders)-1]

	return true, nil
}

func (repo *ProductMemoryRepository) RegisterOrder(userID uint32, order []*OrderProduct) (uint32, error) {
	for _, pr := range order {
		dOrd := &DeliveryOrder{
			UserID:        userID,
			ProductName:   pr.Title,
			Price:         pr.Price,
			DepartureDate: time.Now(),
			LandingDate:   time.Now().Add(4 * 24 * time.Hour),
			Status:        "в пути",
		}
		repo.deliveryOrders = append(repo.deliveryOrders, dOrd)
	}

	return userID, nil
}

func (repo *ProductMemoryRepository) GetDeliveryOrdersByID(userID uint32) ([]*DeliveryOrder, error) {
	orders := []*DeliveryOrder{}
	for _, ord := range repo.deliveryOrders {
		if ord.UserID == userID {
			orders = append(orders, ord)
		}
	}
	return orders, nil
}

func (repo *ProductMemoryRepository) GetRelated(typ string, limit int) ([]*Product, error) {
	prd := []*Product{}
	sort.Slice(repo.data, func(i, j int) bool {
		return repo.data[i].Views > repo.data[j].Views
	})

	for _, item := range repo.data {
		if item.Type == typ {
			prd = append(prd, item)
		}
		if len(prd) == limit {
			break
		}
	}

	return prd, nil
}
