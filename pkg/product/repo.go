package product

import (
	"errors"
	"sort"
	"time"
)

type ProductMemoryRepository struct {
	data    []*Product
	baskets map[uint32]*Basket
	orders  map[uint32]*History
	lastID  uint32
}

func NewMemoryRepo() *ProductMemoryRepository {
	return &ProductMemoryRepository{
		data: []*Product{
			{
				ID:           1,
				Title:        "Хайлайтер для губ",
				Price:        42.33,
				Description:  "Описание1",
				FactoryCount: 5,
				Type:         "косметика",
				Tag:          "Хит",
				ImageURL:     "/static/images/Highlighter.png",
			},
			{
				ID:           2,
				Title:        "Подводка",
				Price:        95.00,
				Description:  "Описание2",
				FactoryCount: 5,
				Type:         "косметика",
				Tag:          "Sale",
				ImageURL:     "/static/images/fit_black.png",
			},
			{
				ID:           3,
				Title:        "Набор из блесков",
				Price:        156.00,
				Description:  "Описание3",
				FactoryCount: 5,
				Type:         "косметика",
				Tag:          "Sale",
				ImageURL:     "/static/images/fit_mini.png",
			},
			{
				ID:           4,
				Title:        "Помада Клубничная",
				Price:        180,
				Description:  "Описание4",
				FactoryCount: 5,
				Type:         "косметика",
				Tag:          "Хит",
				ImageURL:     "/static/images/lip.png",
			},
			{
				ID:           5,
				Title:        "Помада карандаш",
				Price:        200,
				Description:  "Описание5",
				FactoryCount: 5,
				Type:         "косметика",
				Tag:          "Хит",
				ImageURL:     "/static/images/lip_array.png",
			},
			{
				ID:           6,
				Title:        "Помада Матовая",
				Price:        249,
				Description:  "Описание6",
				FactoryCount: 5,
				Type:         "косметика",
				Tag:          "Хит",
				ImageURL:     "/static/images/lip_pink.png",
			},
			{
				ID:           7,
				Title:        "Ингаляторы Нирдош",
				Price:        50000000,
				Description:  "Описание7",
				FactoryCount: 5,
				Type:         "косметика",
				Tag:          "Хит",
				ImageURL:     "/static/images/ing.png",
			},
			{
				ID:           8,
				Title:        "Паста для шугаринга",
				Price:        600,
				Description:  "Описание8",
				FactoryCount: 5,
				Type:         "косметика",
				Tag:          "Хит",
				ImageURL:     "/static/images/Paste.png",
			},
		},
		baskets: map[uint32]*Basket{},
		orders:  map[uint32]*History{},
		lastID:  8,
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

func (repo *ProductMemoryRepository) GetBasketByID(id uint32) (*Basket, error) {
	bsk, ok := repo.baskets[id]
	if !ok {
		return &Basket{}, ErrNoBasket
	}
	return bsk, nil
}

func (repo *ProductMemoryRepository) AddProductToBasket(id uint32, product *Product) (uint32, error) {
	for _, bsk := range repo.baskets {
		if bsk.UserID == id {
			for _, item := range bsk.Products {
				if item.ID == product.ID {
					return 0, nil
				}
			}
			bsk.Products = append(bsk.Products, product)
			bsk.TotalPrice += product.Price
			bsk.TotalCount++
		}
	}
	return product.ID, nil
}

func (repo *ProductMemoryRepository) DeleteProductFromBasket(id, prodID uint32) (bool, error) {
	_, ok := repo.baskets[id]
	if !ok {
		return false, ErrNoBasket
	}
	i := -1
	for idx, item := range repo.baskets[id].Products {
		if item.ID != prodID {
			continue
		}
		i = idx
	}
	if i < 0 {
		return false, nil
	}

	if i < len(repo.baskets[id].Products)-1 {
		copy(repo.baskets[id].Products[i:], repo.baskets[id].Products[i+1:])
	}
	repo.baskets[id].TotalCount = repo.baskets[id].TotalCount - 1
	repo.baskets[id].TotalPrice -= repo.baskets[id].Products[i].Price
	repo.baskets[id].Products[len(repo.baskets[id].Products)-1] = &Product{} // or the zero value of T
	repo.baskets[id].Products = repo.baskets[id].Products[:len(repo.baskets[id].Products)-1]

	return true, nil
}

func (repo *ProductMemoryRepository) AddBasket(id uint32) (uint32, error) {
	bsk := &Basket{
		UserID:   id,
		Products: []*Product{},
	}
	_, ok := repo.baskets[id]
	if !ok {
		repo.baskets[id] = bsk
	}
	return id, nil
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

func (repo *ProductMemoryRepository) RegisterOrder(userID uint32, products []*Product) (uint32, error) {
	_, ok := repo.orders[userID]
	if !ok {
		repo.orders[userID] = &History{
			UserID:   userID,
			Landings: []*LandingInfo{},
		}
	}

	for _, pr := range products {
		lnd := &LandingInfo{
			ProductName:   pr.Title,
			Price:         pr.Price,
			DepartureDate: time.Now(),
			LandingDate:   time.Now().Add(4 * 24 * time.Hour),
			Status:        "в пути",
		}
		repo.orders[userID].Landings = append(repo.orders[userID].Landings, lnd)
	}

	return userID, nil
}

func (repo *ProductMemoryRepository) GetOrders(userID uint32) ([]*LandingInfo, error) {
	ord, ok := repo.orders[userID]
	if !ok {
		return []*LandingInfo{}, nil
	}
	return ord.Landings, nil
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
