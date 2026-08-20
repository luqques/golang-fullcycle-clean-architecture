package graphql

import (
	"context"
	"net/http"

	graphqlgo "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/lucas/clean-orders-challenge/internal/usecase"
)

const Schema = `
schema {
  query: Query
  mutation: Mutation
}

type Query {
  listOrders: [Order!]!
}

type Mutation {
  createOrder(input: CreateOrderInput!): Order!
}

input CreateOrderInput {
  price: Float!
  tax: Float!
}

type Order {
  id: ID!
  price: Float!
  tax: Float!
  finalPrice: Float!
}
`

type Resolver struct {
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewHandler(createUC *usecase.CreateOrderUseCase, listUC *usecase.ListOrdersUseCase) http.Handler {
	resolver := &Resolver{CreateOrderUseCase: createUC, ListOrdersUseCase: listUC}
	schema := graphqlgo.MustParseSchema(Schema, resolver)
	return &relay.Handler{Schema: schema}
}

func (r *Resolver) ListOrders(ctx context.Context) ([]*OrderResolver, error) {
	orders, err := r.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	result := make([]*OrderResolver, 0, len(orders))
	for _, order := range orders {
		result = append(result, &OrderResolver{order: order})
	}

	return result, nil
}

type createOrderArgs struct {
	Input struct {
		Price float64
		Tax   float64
	}
}

func (r *Resolver) CreateOrder(ctx context.Context, args createOrderArgs) (*OrderResolver, error) {
	order, err := r.CreateOrderUseCase.Execute(usecase.CreateOrderInputDTO{
		Price: args.Input.Price,
		Tax:   args.Input.Tax,
	})
	if err != nil {
		return nil, err
	}

	return &OrderResolver{order: order}, nil
}

type OrderResolver struct {
	order usecase.OrderOutputDTO
}

func (r *OrderResolver) ID() graphqlgo.ID { return graphqlgo.ID(r.order.ID) }
func (r *OrderResolver) Price() float64   { return r.order.Price }
func (r *OrderResolver) Tax() float64     { return r.order.Tax }
func (r *OrderResolver) FinalPrice() float64 {
	return r.order.FinalPrice
}
