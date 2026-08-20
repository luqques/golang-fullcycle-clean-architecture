package database

import (
	"database/sql"

	"github.com/lucas/clean-orders-challenge/internal/entity"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	_, err := r.DB.Exec(
		"INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)",
		order.ID,
		order.Price,
		order.Tax,
		order.FinalPrice,
	)
	return err
}

func (r *OrderRepository) List() ([]entity.Order, error) {
	rows, err := r.DB.Query("SELECT id, price, tax, final_price FROM orders ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
