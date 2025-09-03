package events

import (
	"context"
	"ride-sharing/shared/messaging"
)

type TripEventPublisher struct {
	rabbit *messaging.RabbitMQ
}

func NewTripEventPublisher(rabbitmq *messaging.RabbitMQ) *TripEventPublisher {
	return &TripEventPublisher{rabbit: rabbitmq}
}

func (t *TripEventPublisher) PublishTripCreated(ctx context.Context) error {
	return t.rabbit.PublishingMessage(ctx, "hallo", "hallo world")
}
