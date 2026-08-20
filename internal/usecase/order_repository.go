package usecase

import "github.com/lucas/clean-orders-challenge/internal/entity"

type OrderRepository interface {
	Save(order *entity.Order) error
	List() ([]entity.Order, error)
}
