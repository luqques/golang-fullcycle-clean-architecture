package usecase

type ListOrdersUseCase struct {
	OrderRepository OrderRepository
}

func NewListOrdersUseCase(repository OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{OrderRepository: repository}
}

func (uc *ListOrdersUseCase) Execute() ([]OrderOutputDTO, error) {
	orders, err := uc.OrderRepository.List()
	if err != nil {
		return nil, err
	}

	output := make([]OrderOutputDTO, 0, len(orders))
	for _, order := range orders {
		output = append(output, OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return output, nil
}
