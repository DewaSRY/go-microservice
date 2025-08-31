package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/events"
	grpcHandler "ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"

	"ride-sharing/services/trip-service/internal/service"
	"syscall"

	grpcServer "google.golang.org/grpc"
)

const GrpcAddr = ":9093"

func main() {
	rabbitmqUri := env.GetString("RABBITMQ_URI", "amqp://guess:guess@rabbitmq:5672/")

	ctx, cancel := context.WithCancel(context.Background())

	log.Println("Starting API Gateway")

	repository := repository.NewInMemoryRepository()
	service := service.NewService(repository)

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

	rabbitmq, err := messaging.NewRabbitMQ(rabbitmqUri)
	if err != nil {
		log.Fatal("failed_to_connect_to_rabbitmq: ")
	}
	defer rabbitmq.Close()

	publisher := events.NewTripEventPublisher(rabbitmq)

	grpcServer := grpcServer.NewServer()
	grpcHandler.NewGRPCHandler(grpcServer, service, publisher)

	log.Printf("Starting gRPC server Trip service on port %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed_to_serve:%v", err)
			cancel()
		}
	}()

	<-ctx.Done()
}
