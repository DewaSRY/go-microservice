package main

import (
	"fmt"
	"log"
	"net/http"
	grpcclient "ride-sharing/services/api-gateway/grpc_client"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"ride-sharing/shared/proto/driver"

	"github.com/gorilla/websocket"
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

func handlerRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connect ")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed :%v", err)
		return
	}

	userId := r.URL.Query().Get("userID")
	if userId == "" {
		log.Println("no_user_id_provided")
		return
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Printf("error_reading_message: %v\n", err)
			break
		}

		log.Printf("received_messages: %s", message)
	}
}

func handleDriverWebScoket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
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

	driverService, err := grpcclient.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		driverService.Client.UnRegisterDriver(ctx, &driver.RegisterDriverRequest{
			DriverId:    userId,
			PackageSlug: packageSlug,
		})
		driverService.Close()
		log.Println("Driver unregistered: ", userId)
	}()

	driverData, err := driverService.Client.RegisterDriver(ctx, &driver.RegisterDriverRequest{
		DriverId:    userId,
		PackageSlug: packageSlug,
	})
	if err != nil {
		log.Printf("Error registering driver: %v", err)
		return
	}

	fmt.Print(driverData)

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Printf("error_reading_message: %v\n", err)
			break
		}

		log.Printf("received_messages: %s", message)
	}
}
