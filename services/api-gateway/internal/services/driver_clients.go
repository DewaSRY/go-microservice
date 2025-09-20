package services

import (
	"fmt"
	"os"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DriverServiceClient struct {
	Client pb.DriverServiceClient
	Conn   *grpc.ClientConn
}

func NewDriverServiceClient() (*DriverServiceClient, error) {
	driverServiceURL := os.Getenv("DRIVER_SERVICE_URL")
	if driverServiceURL == "" {
		driverServiceURL = "driver-service:9092"
	}

	// dialOptions := append(
	// 	tracing.DialOptionsWithTracing(),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )

	conn, err := grpc.NewClient(driverServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed_create_new_client :%v", err)
	}

	client := pb.NewDriverServiceClient(conn)
	return &DriverServiceClient{
		Conn:   conn,
		Client: client,
	}, nil
}

func (c *DriverServiceClient) Close() {
	if c.Client != nil {
		if err := c.Conn.Close(); err != nil {
			return
		}
	}
}
