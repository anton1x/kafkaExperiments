package usecase

import (
	"context"
	"log/slog"

	"kafkaExperiments/internal/domain/order/entities"
	"kafkaExperiments/internal/domain/order/events"
	"kafkaExperiments/internal/order"
)

type Order struct {
	repo     order.Repo
	logger   *slog.Logger
	eventBus order.EventBus
}

func NewOrderUseCase(logger *slog.Logger, repo order.Repo, bus order.EventBus) *Order {
	return &Order{
		repo:     repo,
		logger:   logger,
		eventBus: bus,
	}
}

func (svc *Order) Create(ctx context.Context, newOrder *entities.NewOrder) (*entities.Order, error) {
	order, err := svc.repo.CreateOrder(ctx, newOrder)

	if err != nil {
		svc.logger.Log(ctx, slog.LevelError, "error creating order")
		return nil, err
	}

	err = svc.eventBus.Publish(ctx, events.NewOrderCreatedEvent(events.OrderCreatedParams{
		Id:         order.Id,
		CustomerId: order.CustomerId,
	}))

	if err != nil {
		svc.logger.Error("error publish", "error", err)
	}

	return order, nil
}
