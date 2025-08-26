package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	grpcHandler "ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"

	"ride-sharing/services/trip-service/internal/service"
	"syscall"

	grpcServer "google.golang.org/grpc"
)

const GrpcAddr = ":9093"

func main() {
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

	grpcServer := grpcServer.NewServer()
	grpcHandler.NewGRPCHandler(grpcServer, service)

	log.Printf("Starting gRPC server Trip service on port %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed_to_serve:%v", err)
			cancel()
		}
	}()

	<-ctx.Done()
}
