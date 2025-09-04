package events

import (
	"context"
	"log"
	"ride-sharing/shared/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type TripEventConnsummer struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripConsummer(rabbitmq *messaging.RabbitMQ) *TripEventConnsummer {
	return &TripEventConnsummer{
		rabbitmq: rabbitmq,
	}
}

func (t *TripEventConnsummer) Listen() error {
	return t.rabbitmq.ConsummeMessages(messaging.FindAvailableDriversQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		log.Printf("deliver received message : %v", msg)
		return nil
	})
}
