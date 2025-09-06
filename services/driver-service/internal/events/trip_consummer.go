package events

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"ride-sharing/services/driver-service/internal/service"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type TripEventConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  *service.Service
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ, service *service.Service) *TripEventConsumer {
	return &TripEventConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (t *TripEventConsumer) Listen() error {
	return t.rabbitmq.ConsumeMessages(messaging.FindAvailableDriversQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var tripEvent contracts.AmqpMessage

		if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
			log.Printf("failed_to_unmarshal_message:%v", err)
			return err
		}

		var payload messaging.TripEventData
		if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
			log.Printf("failed_to_unmarshal_payload:%v", err)
		}

		switch msg.RoutingKey {
		case contracts.TripEventCreated, contracts.TripEventDriverNotInterested:
			return t.handleFindAndNotifyDrivers(ctx, payload)
		}

		log.Printf("deliver_received_message : %v", payload)
		return nil
	})
}

func (t *TripEventConsumer) handleFindAndNotifyDrivers(ctx context.Context, payload messaging.TripEventData) error {
	suitableIDs := t.service.FindAvailableDrivers(payload.Trip.SelectedFare.PackageSlug)

	log.Printf("Found suitable drivers %v", len(suitableIDs))

	if len(suitableIDs) == 0 {
		// Notify the driver that no drivers are available
		if err := t.rabbitmq.PublishingMessage(ctx, contracts.TripEventNoDriversFound, contracts.AmqpMessage{
			OwnerID: payload.Trip.UserID,
		}); err != nil {
			log.Printf("Failed to publish message to exchange: %v", err)
			return err
		}

		return nil
	}

	// suitableDriverID := suitableIDs[0]
	randomIndex := rand.Intn(len(suitableIDs))
	suitableDriverID := suitableIDs[randomIndex]

	marshalledEvent, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Notify the driver about a potential trip
	if err := t.rabbitmq.PublishingMessage(ctx, contracts.DriverCmdTripRequest, contracts.AmqpMessage{
		OwnerID: suitableDriverID,
		Data:    marshalledEvent,
	}); err != nil {
		log.Printf("Failed to publish message to exchange: %v", err)
		return err
	}

	return nil
}
