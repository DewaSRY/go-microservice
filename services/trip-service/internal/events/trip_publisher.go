package events

import (
	"context"
	"encoding/json"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/mapper"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/types"
)

type TripEventPublisher struct {
	rabbit *messaging.RabbitMQ
}

func NewTripEventPublisher(rabbitmq *messaging.RabbitMQ) *TripEventPublisher {
	return &TripEventPublisher{rabbit: rabbitmq}
}

func (t *TripEventPublisher) PublishTripCreated(ctx context.Context, trip types.TripModel) error {
	payload := messaging.TripEventData{
		Trip: mapper.MappedTripModelToProtoTripModel(&trip),
	}

	tripEventJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return t.rabbit.PublishingMessage(ctx, contracts.TripEventCreated, contracts.AmqpMessage{
		OwnerID: trip.UserId,
		Data:    tripEventJSON,
	})
}
