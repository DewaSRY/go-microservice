package messaging

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

type MessageHandler func(context.Context, amqp.Delivery) error

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
		true,    //durable
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

func (r *RabbitMQ) ConsummeMessages(queueName string, handler MessageHandler) error {

	// Set prefetch count to 1 for fair dispatch
	// This tells RabbitMQ not to give more than one message to a service at a time.
	// The worker will only get the next message after it has acknowledged the previous one.
	err := r.Channel.Qos(
		1,     // prefetchCount: Limit to 1 unacknowledged message per consumer
		0,     // prefetchSize: No specific limit on message size
		false, // global: Apply prefetchCount to each consumer individually
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	msgs, err := r.Channel.Consume(
		"hallo", // quue
		"",      // consummer
		false,   //auto-ack
		false,   //exclusive
		false,   //no-local
		false,   //no-await
		nil,     //args
	)

	if err != nil {
		return err
	}
	ctx := context.Background()

	go func() {
		for msg := range msgs {
			log.Printf("receive_message :%s", msg.Body)

			if err := handler(ctx, msg); err != nil {
				log.Printf("error_failed_to_handle_message:%v", err)
				//Neck the message. set requeue to false to avoid immediate redelivery loops.
				//consider a dead-letter exchange (DLQ) or a more sophisticated retry mechanism for production
				if neckErr := msg.Nack(false, false); neckErr != nil {
					log.Printf("error_failed_tonack_message:%v", neckErr)
				}
				continue
			}

			//only acknowledgment if the handler is success
			if ackErr := msg.Ack(false); ackErr != nil {
				log.Printf("error_failed_to_ack_message :%v", ackErr)
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}

	if r.Channel != nil {
		r.Channel.Close()
	}
}
