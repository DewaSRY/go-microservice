package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"

	"ride-sharing/shared/mapper"

	"google.golang.org/grpc"
)

type GRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService) *GRPCHandler {
	handler := &GRPCHandler{
		service: service,
	}
	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *GRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := mapper.MappedLocationToCoridinate(req.GetStartLocation())
	destination := mapper.MappedLocationToCoridinate(req.GetEndLocation())

	trip_route, err := h.service.GetRoute(ctx, pickup, destination)

	if err != nil {
		log.Println(err)
	}

	routes := make([]*pb.Route, len(trip_route.Routes))

	for i, r := range trip_route.Routes {
		routes[i] = mapper.MappedRouteToProtoroute(&r)
	}

	if err != nil {
		log.Println(err)
	}

	return &pb.PreviewTripResponse{
		Route:     routes,
		RideFares: []*pb.RideFare{},
	}, nil
}
