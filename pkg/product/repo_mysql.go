package product

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	views     = "views"
	createdAt = "created_at"
)

type ProductMysqlRepository struct {
	DB *sql.DB
}

func NewMysqlRepo(db *sql.DB) *ProductMysqlRepository {
	return &ProductMysqlRepository{DB: db}
}

func (repo *ProductMysqlRepository) GetAll(orderBy string) ([]*Product, error) {
	query := "SELECT id, title, price, tag, type, description, count, created_at, views, image_url FROM products"
	switch orderBy {
	case views:
		query = "SELECT id, title, price, tag, type, description, count, created_at, views, image_url FROM products ORDER BY views DESC"
	case createdAt:
		query = "SELECT id, title, price, tag, type, description, count, created_at, views, image_url FROM products ORDER BY created_at DESC"
	}
	products := []*Product{}
	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		prd := &Product{}
		err = rows.Scan(&prd.ID, &prd.Title, &prd.Price, &prd.Tag, &prd.Type, &prd.Description, &prd.Count, &prd.CreatedAt, &prd.Views, &prd.ImageURL)
		if err != nil {
			return nil, err
		}
		products = append(products, prd)
	}
	return products, nil
}

func (repo *ProductMysqlRepository) GetByID(id uint32) (*Product, error) {
	prd := &Product{}
	// QueryRow сам закрывает коннект
	err := repo.DB.
		QueryRow("SELECT id, title, price, tag, type, description, count, created_at, views, image_url FROM products WHERE id = ?", id).
		Scan(&prd.ID, &prd.Title, &prd.Price, &prd.Tag, &prd.Type, &prd.Description, &prd.Count, &prd.CreatedAt, &prd.Views, &prd.ImageURL)
	if err != nil {
		return nil, err
	}
	return prd, nil
}

func (repo *ProductMysqlRepository) Add(product *Product) (uint32, error) {
	result, err := repo.DB.Exec(
		"INSERT INTO products (`title`, `price`, `tag`, `type`, `description`, `count`, `created_at`, `views`, `image_url`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		product.Title,
		product.Price,
		product.Tag,
		product.Type,
		product.Description,
		product.Count,
		product.CreatedAt,
		product.Views,
		product.ImageURL,
	)
	if err != nil {
		fmt.Print(err.Error())
		return 0, err
	}
	res, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(res), nil
}

func (repo *ProductMysqlRepository) Update(newProduct *Product) (bool, error) {
	_, err := repo.DB.Exec(
		"UPDATE products SET"+
			"`title` = ?"+
			",`price` = ?"+
			",`tag` = ?"+
			",`type` = ?"+
			",`description` = ?"+
			",`count` = ?"+
			",`created_at` = ?"+
			",`views` = ?"+
			",`image_url` = ?"+
			"WHERE id = ?",
		newProduct.Title,
		newProduct.Price,
		newProduct.Tag,
		newProduct.Type,
		newProduct.Description,
		newProduct.Count,
		newProduct.CreatedAt,
		newProduct.Views,
		newProduct.ImageURL,
		newProduct.ID,
	)
	if err != nil {
		fmt.Print(err.Error())
		return false, err
	}

	return true, nil
	//return result.RowsAffected()
}

func (repo *ProductMysqlRepository) Delete(id uint32) (bool, error) {
	_, err := repo.DB.Exec(
		"DELETE FROM products WHERE id = ?",
		id,
	)
	if err != nil {
		return false, err
	}
	return true, nil
	//return result.RowsAffected()
}

func (repo *ProductMysqlRepository) GetOrdersByID(id uint32) (*Order, error) {
	orders := &Order{}
	rows, err := repo.DB.Query("SELECT user_id, order_id, title, description, price, image_url FROM orders WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		ord := &OrderProduct{}
		err = rows.Scan(&ord.UserID, &ord.ID, &ord.Title, &ord.Description, &ord.Price, &ord.ImageURL)
		if err != nil {
			return nil, err
		}
		orders.Products = append(orders.Products, ord)
		orders.TotalPrice += ord.Price
	}
	return orders, nil
}

func (repo *ProductMysqlRepository) AddOrder(id uint32, product *Product) (uint32, error) {
	result, err := repo.DB.Exec(
		"INSERT IGNORE INTO orders (`user_id`, `order_id`, `title`, `description`, `price`, `image_url`) VALUES (?, ?, ?, ?, ?, ?)",
		id,
		product.ID,
		product.Title,
		product.Description,
		product.Price,
		product.ImageURL,
	)
	if err != nil {
		return 0, err
	}
	res, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return uint32(res), nil
}

func (repo *ProductMysqlRepository) DeleteOrder(id, prodID uint32) (bool, error) {
	_, err := repo.DB.Exec(
		"DELETE FROM orders WHERE user_id = ? AND order_id = ?",
		id, prodID,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *ProductMysqlRepository) RegisterOrder(userID uint32, order []*OrderProduct) (uint32, error) {
	deliveryOrders := []*DeliveryOrder{}
	for _, ord := range order {
		dOrd := &DeliveryOrder{
			UserID:        userID,
			ProductName:   ord.Title,
			Price:         ord.Price,
			DepartureDate: time.Now(),
			LandingDate:   time.Now().Add(4 * 24 * time.Hour),
			Status:        "в пути",
		}
		deliveryOrders = append(deliveryOrders, dOrd)
	}

	sqlStr := "INSERT INTO deliveryOrders(`user_id`, `product_name`, `price`, `departure_date`, `landing_date`, `status`) VALUES "
	vals := []interface{}{}

	for _, row := range deliveryOrders {
		sqlStr += "(?, ?, ?, ?, ?, ?),"
		vals = append(vals, row.UserID, row.ProductName, row.Price, row.DepartureDate, row.LandingDate, row.Status)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	stmt, err := repo.DB.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}

	_, err = stmt.Exec(vals...)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (repo *ProductMysqlRepository) GetDeliveryOrdersByID(userID uint32) ([]*DeliveryOrder, error) {
	orders := []*DeliveryOrder{}
	rows, err := repo.DB.Query("SELECT user_id, product_name, price, departure_date, landing_date, status FROM deliveryOrders WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		ord := &DeliveryOrder{}
		err = rows.Scan(&ord.UserID, &ord.ProductName, &ord.Price, &ord.DepartureDate, &ord.LandingDate, &ord.Status)
		if err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	return orders, nil
}

func (repo *ProductMysqlRepository) GetRelated(typ string, limit int) ([]*Product, error) {
	query := fmt.Sprintf("SELECT id, title, price, tag, type, description, count, created_at, views, image_url FROM products WHERE type = ? LIMIT %d", limit)
	products := []*Product{}
	rows, err := repo.DB.Query(query, typ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		prd := &Product{}
		err = rows.Scan(&prd.ID, &prd.Title, &prd.Price, &prd.Tag, &prd.Type, &prd.Description, &prd.Count, &prd.CreatedAt, &prd.Views, &prd.ImageURL)
		if err != nil {
			return nil, err
		}
		products = append(products, prd)
	}
	return products, nil
}
