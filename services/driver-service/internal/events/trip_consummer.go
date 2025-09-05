package events

import (
	"context"
	"encoding/json"
	"log"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type TripEventConsumer struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ) *TripEventConsumer {
	return &TripEventConsumer{
		rabbitmq: rabbitmq,
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

		log.Printf("deliver_received_message : %v", payload)
		return nil
	})
}
