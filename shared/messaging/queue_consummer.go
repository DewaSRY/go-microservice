package messaging

import (
	"encoding/json"
	"log"

	"ride-sharing/shared/contracts"
	"ride-sharing/shared/ws"
)

type QueueConsumer interface {
	Start() error
}

type queueConsumer struct {
	rb        RabbitMQClient
	connMgr   ws.ConnectionManager
	queueName string
}

func NewQueueConsumer(rb RabbitMQClient, connMgr ws.ConnectionManager, queueName string) QueueConsumer {
	return &queueConsumer{
		rb:        rb,
		connMgr:   connMgr,
		queueName: queueName,
	}
}

func (qc *queueConsumer) Start() error {
	msgs, err := qc.rb.CreateConsumer(
		qc.queueName,
		true,
		false,
		false,
		false,
	)

	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var msgBody contracts.AmqpMessage
			if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
				log.Println("Failed to unmarshal message:", err)
				continue
			}

			userID := msgBody.OwnerID

			var payload any
			if msgBody.Data != nil {
				if err := json.Unmarshal(msgBody.Data, &payload); err != nil {
					log.Println("Failed to unmarshal payload:", err)
					continue
				}
			}

			clientMsg := contracts.WSMessage{
				Type: msg.RoutingKey,
				Data: payload,
			}

			if err := qc.connMgr.Emit(userID, clientMsg); err != nil {
				log.Printf("Failed to send message to user %s: %v", userID, err)
			}
		}
	}()

	return nil
}
