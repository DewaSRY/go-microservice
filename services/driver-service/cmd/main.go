package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"ride-sharing/services/driver-service/internal/events"
	"ride-sharing/services/driver-service/internal/infrastructure/grpc"
	"ride-sharing/services/driver-service/internal/service"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"

	grpcServer "google.golang.org/grpc"
)

const GrpcAddr = ":9092"

func main() {
	rabbitmqUri := env.GetString("RABBITMQ_URI", "amqp://guess:guess@rabbitmq:5672/")
	ctx, cancel := context.WithCancel(context.Background())
	driverService := service.NewService()
	log.Println("Starting API Gateway")
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed_to_listen : %v", err)
	}

	rabbit, err := messaging.NewRabbitMQ(rabbitmqUri)
	if err != nil {
		log.Fatal("failed_to_connect_to_rabbitmq: ")
	}
	defer rabbit.Close()

	tripConnnsummer := events.NewTripConsummer(rabbit)

	go func() {
		fmt.Print("statt to listen")

		if err := tripConnnsummer.Listen(); err != nil {
			log.Fatalf("failed_to_connsume_the_queue :%v", err)
		}
	}()

	//Create Grpc service
	grpcServer := grpcServer.NewServer()
	grpc.NewGrpcHandler(grpcServer, driverService)

	log.Printf("Starting gRPC server Driver service on port %s", lis.Addr().String())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed_to_serve:%v", err)
			cancel()
		}
	}()

	<-ctx.Done()
}
