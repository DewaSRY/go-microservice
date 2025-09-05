package grpc

import (
	"context"
	"log"
	"ride-sharing/services/driver-service/internal/service"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	pb.UnimplementedDriverServiceServer
	Service *service.Service
}

func NewGrpcHandler(s *grpc.Server, service *service.Service) {
	handler := &GrpcHandler{
		Service: service,
	}

	pb.RegisterDriverServiceServer(s, handler)
}

func (h *GrpcHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	log.Printf("register_driver: %v", req)

	driver, err := h.Service.RegisterDriver(req.GetDriverId(), req.GetPackageSlug())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register driver")
	}

	return &pb.RegisterDriverResponse{
		Driver: driver,
	}, nil
}

func (h *GrpcHandler) UnRegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	log.Printf("unregister_driver: %v", req)
	h.Service.UnregisterDriver(req.GetDriverId())

	return &pb.RegisterDriverResponse{
		Driver: &pb.Driver{
			Id: req.GetDriverId(),
		},
	}, nil
}
