package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/shared/contracts"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient interface {
	PublishingMessage(ctx context.Context, routingKey string, message contracts.AmqpMessage) error
	ConsumeMessages(queueName string, handler MessageHandler) error
	CreateConsumer(queueName string, autoAck bool, exclusive bool, noLocal bool, noWait bool) (<-chan amqp.Delivery, error)
	Close()
}

type rabbitMQClientImpl struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

const (
	TripExchange = "trip"
)

type MessageHandler func(context.Context, amqp.Delivery) error

func NewRabbitMQ(urlString string) (RabbitMQClient, error) {
	conn, err := amqp.Dial(urlString)
	if err != nil {
		return nil, fmt.Errorf("failed_connect_to_rabbitmq")
	}

	ch, err := conn.Channel()

	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed_to_create_channel")
	}

	rmq := &rabbitMQClientImpl{
		conn:    conn,
		Channel: ch,
	}

	if err := rmq.setupExchangesAndQueues(); err != nil {
		rmq.Close()
		return nil, fmt.Errorf("failed_to_setup_and_exchange_queue")
	}

	return rmq, nil
}

func NewRabbitMQ2(urlString string) (*rabbitMQClientImpl, error) {
	conn, err := amqp.Dial(urlString)
	if err != nil {
		return nil, fmt.Errorf("failed_connect_to_rabbitmq")
	}

	ch, err := conn.Channel()

	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed_to_create_channel")
	}

	rmq := &rabbitMQClientImpl{
		conn:    conn,
		Channel: ch,
	}

	if err := rmq.setupExchangesAndQueues(); err != nil {
		rmq.Close()
		return nil, fmt.Errorf("failed_to_setup_and_exchange_queue")
	}

	return rmq, nil
}

func (r *rabbitMQClientImpl) setupExchangesAndQueues() error {
	err := r.Channel.ExchangeDeclare(
		TripExchange, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return fmt.Errorf("failed to declare exchange: %s: %v", TripExchange, err)
	}

	if err := r.declareAndBindQueue(
		FindAvailableDriversQueue,
		[]string{
			contracts.TripEventCreated, contracts.TripEventDriverNotInterested,
		},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		DriverCmdTripRequestQueue,
		[]string{contracts.DriverCmdTripRequest},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		DriverTripResponseQueue,
		[]string{contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		NotifyDriverNoDriversFoundQueue,
		[]string{contracts.TripEventNoDriversFound},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		NotifyDriverAssignQueue,
		[]string{contracts.TripEventDriverAssigned},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		PaymentTripResponseQueue,
		[]string{contracts.PaymentCmdCreateSession},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		NotifyPaymentSessionCreatedQueue,
		[]string{contracts.PaymentEventSessionCreated},
		TripExchange,
	); err != nil {
		return err
	}

	return nil
}

func (r *rabbitMQClientImpl) PublishingMessage(ctx context.Context, routingKey string, message contracts.AmqpMessage) error {

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	return r.Channel.PublishWithContext(
		ctx,
		TripExchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonMsg,
		})
}

func (r *rabbitMQClientImpl) CreateConsumer(queueName string, autoAck bool, exclusive bool, noLocal bool, noWait bool) (<-chan amqp.Delivery, error) {
	return r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     //auto-ack
		false,     //exclusive
		false,     //no-local
		false,     //no-await
		nil,       //args
	)
}

func (r *rabbitMQClientImpl) ConsumeMessages(queueName string, handler MessageHandler) error {
	// Set prefetch count to 1 for fair dispatch
	// This tells RabbitMQ not to give more than one message to a service at a time.
	// The worker will only get the next message after it has acknowledged the previous one.
	err := r.Channel.Qos(
		1,     // prefetchCount: Limit to 1 unacknowledged message per consumer
		0,     // prefetchSize: No specific limit on message size
		false, // global: Apply prefetchCount to each consumer individually
	)
	if err != nil {
		return fmt.Errorf("failed_to_set_QoS: %v", err)
	}

	msgs, err := r.CreateConsumer(
		queueName, // queue
		false,     //auto-ack
		false,     //exclusive
		false,     //no-local
		false,     //no-await
	)

	if err != nil {
		log.Fatalf("failed_to_create_consumer:%v", err)
		return err
	}

	ctx := context.Background()

	go func() {
		for msg := range msgs {
			if err := handler(ctx, msg); err != nil {
				log.Printf("error_failed_to_handle_message:%v", err)
				//Neck the message. set requeue to false to avoid immediate redelivery loops.
				//consider a dead-letter exchange (DLQ) or a more sophisticated retry mechanism for production
				if neckErr := msg.Nack(false, false); neckErr != nil {
					log.Printf("error_failed_to_nack_message:%v", neckErr)
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

func (r *rabbitMQClientImpl) declareAndBindQueue(queueName string, messageTypes []string, exchange string) error {
	q, err := r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range messageTypes {
		if err := r.Channel.QueueBind(
			q.Name,   // queue name
			msg,      // routing key
			exchange, // exchange
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed_to_bind_queue_to %s: %v", queueName, err)
		}
	}

	return nil
}

func (r *rabbitMQClientImpl) Close() {
	if r.conn != nil {
		r.conn.Close()
	}

	if r.Channel != nil {
		r.Channel.Close()
	}
}
