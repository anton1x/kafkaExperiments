package order

import (
	"context"
	"kafkaExperiments/internal/domain/order/entities"
	"kafkaExperiments/internal/infrastructure/eventBus"
)

type Repo interface {
	CreateOrder(ctx context.Context, order *entities.NewOrder) (*entities.Order, error)
}

type EventBus interface {
	Publish(ctx context.Context, msg eventBus.EventMessage) error
}
