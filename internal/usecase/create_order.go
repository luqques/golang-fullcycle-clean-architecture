package usecase

import "github.com/lucas/clean-orders-challenge/internal/entity"

type CreateOrderInputDTO struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository OrderRepository
}

func NewCreateOrderUseCase(repository OrderRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{OrderRepository: repository}
}

func (uc *CreateOrderUseCase) Execute(input CreateOrderInputDTO) (OrderOutputDTO, error) {
	order, err := entity.NewOrder(input.Price, input.Tax)
	if err != nil {
		return OrderOutputDTO{}, err
	}

	if err := uc.OrderRepository.Save(order); err != nil {
		return OrderOutputDTO{}, err
	}

	return OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
