package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	grpcClient "ride-sharing/services/api-gateway/grpc_client"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"ride-sharing/shared/messaging"
	"ride-sharing/shared/proto/driver"

	"github.com/gorilla/websocket"
)

var (
	connManager = messaging.NewConnectionManager()
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Driver struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profilePicture"`
	CarPlate       string `json:"carPlage"`
	PackageSlug    string `json:"packageSlug"`
}

type WsHandler struct {
	rabbitMq *messaging.RabbitMQ
}

func NewWsHandler(rabbitmq *messaging.RabbitMQ) *WsHandler {
	return &WsHandler{rabbitMq: rabbitmq}
}

func (t *WsHandler) HandlerRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connect ")
	conn, err := connManager.Upgrade(w, r)

	if err != nil {
		log.Printf("websocket_upgrade_failed :%v", err)
		return
	}

	userId := r.URL.Query().Get("userID")
	if userId == "" {
		log.Println("no_user_id_provided")
		return
	}

	defer conn.Close()

	// Add connection to manager
	connManager.Add(userId, conn)
	defer connManager.Remove(userId)

	// Initialize queue consumers
	queues := []string{
		messaging.NotifyDriverNoDriversFoundQueue,
		messaging.NotifyDriverAssignQueue,
	}

	for _, q := range queues {
		consumer := messaging.NewQueueConsumer(t.rabbitMq, connManager, q)

		if err := consumer.Start(); err != nil {
			log.Printf("Failed to start consumer for queue: %s: err: %v", q, err)
		}
	}

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Printf("error_reading_message: %v\n", err)
			break
		}

		log.Printf("received_messages: %s", message)
	}
}

func (t *WsHandler) HandleDriverWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := connManager.Upgrade(w, r)

	if err != nil {
		log.Printf("WebSocket upgrade failed :%v", err)
		return
	}

	userId := r.URL.Query().Get("userID")
	if userId == "" {
		log.Println("no_user_id_provided")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Println("no_package_slug_provided")
		return
	}

	// Add connection to manager
	connManager.Add(userId, conn)

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			Id:             userId,
			Name:           "Tiago",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate:       "hallo",
			PackageSlug:    packageSlug,
		},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("error_write_the_message :%v\n", err)
		return
	}

	defer conn.Close()
	ctx := r.Context()

	driverService, err := grpcClient.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		connManager.Remove(userId)
		driverService.Client.UnRegisterDriver(ctx, &driver.RegisterDriverRequest{
			DriverId:    userId,
			PackageSlug: packageSlug,
		})
		driverService.Close()
		log.Println("driver_unregister: ", userId)
	}()

	driverData, err := driverService.Client.RegisterDriver(ctx, &driver.RegisterDriverRequest{
		DriverId:    userId,
		PackageSlug: packageSlug,
	})
	if err != nil {
		log.Printf("error_register_driver: %v", err)
		return
	}

	if err := connManager.SendMessage(userId, contracts.WSMessage{
		Type: contracts.DriverCmdRegister,
		Data: driverData.Driver,
	}); err != nil {
		log.Printf("error_sending_message: %v", err)
		return
	}

	queues := []string{
		messaging.DriverCmdTripRequestQueue,
	}

	for _, q := range queues {
		consumer := messaging.NewQueueConsumer(t.rabbitMq, connManager, q)

		if err := consumer.Start(); err != nil {
			log.Printf("failed_to_start_consumer_for_queue: %s: err: %v", q, err)
		}
	}

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Printf("error_reading_message: %v\n", err)
			break
		}

		type driverMessage struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}

		var driverMsg driverMessage
		if err := json.Unmarshal(message, &driverMsg); err != nil {
			log.Printf("error_unmarshaling_driver_message: %v", err)
			continue
		}

		switch driverMsg.Type {
		case contracts.DriverCmdLocation:
			continue
		case contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline:
			if err := t.rabbitMq.PublishingMessage(ctx, driverMsg.Type, contracts.AmqpMessage{
				OwnerID: userId,
				Data:    driverMsg.Data,
			}); err != nil {
				log.Printf("Error_publishing_message_to_rabbitMQ: %v", err)
			}
		default:
			log.Printf("Unknown_message_type: %s", driverMsg.Type)
		}

		log.Printf("received_messages: %s", message)
	}
}
