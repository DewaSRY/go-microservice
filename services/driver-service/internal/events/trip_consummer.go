package events

import (
	"context"
	"log"
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
		log.Printf("deliver_received_message : %v", msg)
		return nil
	})
}
