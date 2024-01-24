package events

import "kafkaExperiments/internal/infrastructure/eventBus"

const (
	eventPrefix = "order."

	OrderCreatedEvent = eventPrefix + "created"
)

type OrderCreatedParams struct {
	Id         string `json:"id"`
	CustomerId string `json:"customerId"`
}

func NewOrderCreatedEvent(params OrderCreatedParams) eventBus.EventMessage {
	return eventBus.EventMessage{
		ID:     OrderCreatedEvent,
		Event:  OrderCreatedEvent,
		Params: params,
	}
}
