package events

import (
	"context"
	"encoding/json"
	"ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
)

type TripEventPublisher struct {
	rabbit messaging.RabbitMQClient
}

func NewTripEventPublisher(rabbitmq messaging.RabbitMQClient) *TripEventPublisher {
	return &TripEventPublisher{rabbit: rabbitmq}
}

func (t *TripEventPublisher) PublishTripCreated(ctx context.Context, trip types.TripModel) error {
	payload := messaging.TripEventData{
		Trip: trip.ToTripProtoTrip(),
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
