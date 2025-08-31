package messaging

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(urlString string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(urlString)
	if err != nil {
		return nil, fmt.Errorf("failed_connect_to_rabbitmq")
	}

	ch, err := conn.Channel()

	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed_to_create_channel")
	}

	rmq := &RabbitMQ{
		conn:    conn,
		Channel: ch,
	}

	if err := rmq.setupExchngesAndQueues(); err != nil {
		rmq.Close()
		return nil, fmt.Errorf("failed_to_setup_and_exchange_queue")
	}

	return rmq, nil
}

func (r *RabbitMQ) setupExchngesAndQueues() error {
	_, err := r.Channel.QueueDeclare(
		"hallo", //name
		false,   //durable
		false,   //delete when unused
		false,   //exclusive
		false,   //no-wait
		nil,     // argument
	)

	if err != nil {
		return fmt.Errorf("failed_to_declare_queue")
	}

	return nil
}

func (r *RabbitMQ) PublishingMessage(ctx context.Context, routingKey string, message string) error {
	return r.Channel.PublishWithContext(
		ctx,
		"",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plane",
			Body:        []byte(message),
		})
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}

	if r.Channel != nil {
		r.Channel.Close()
	}
}
