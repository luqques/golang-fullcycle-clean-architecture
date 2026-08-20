package grpc

import (
	"context"

	"github.com/lucas/clean-orders-challenge/internal/usecase"
	pb "github.com/lucas/clean-orders-challenge/proto"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	ListOrdersUseCase *usecase.ListOrdersUseCase
}

func NewOrderService(listUC *usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{ListOrdersUseCase: listUC}
}

func (s *OrderService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	response := &pb.ListOrdersResponse{Orders: make([]*pb.Order, 0, len(orders))}
	for _, order := range orders {
		response.Orders = append(response.Orders, &pb.Order{
			Id:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return response, nil
}
