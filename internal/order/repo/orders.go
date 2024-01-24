package repo

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"kafkaExperiments/internal/domain/order/entities"
)

type OrdersMockRepo struct {
}

func (r OrdersMockRepo) CreateOrder(ctx context.Context, order *entities.NewOrder) (*entities.Order, error) {
	id, _ := uuid.GenerateUUID()
	return &entities.Order{
		NewOrder: entities.NewOrder{
			CustomerId: "aaa",
		},
		Id: id,
	}, nil
}
