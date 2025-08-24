package grpcclient

import (
	"fmt"
	"os"
	pb "ride-sharing/shared/proto/trip"

	"google.golang.org/grpc"
)

type TripServiceClient struct {
	Client pb.TripServiceClient
	Conn   *grpc.ClientConn
}

func NewTripServiceClient() (*TripServiceClient, error) {
	tripServiceURL := os.Getenv("TRIP_SERVICE-URL")
	if tripServiceURL == "" {
		tripServiceURL = "trip-service:9093"
	}
	conn, err := grpc.NewClient(tripServiceURL)
	if err != nil {
		return nil, fmt.Errorf("failed_create_new_client :%v", err)
	}

	client := pb.NewTripServiceClient(conn)
	return &TripServiceClient{
		Conn:   conn,
		Client: client,
	}, nil
}

func (c *TripServiceClient) Close() {
	if c.Client != nil {
		if err := c.Conn.Close(); err != nil {
			return
		}
	}
}
