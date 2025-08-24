package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"dario.cat/mergo"
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
	pickup := mappedLocationToCoridinate(req.GetStartLocation())
	destination := mappedLocationToCoridinate(req.GetEndLocation())

	trip_route, err := h.service.GetRoute(ctx, pickup, destination)

	if err != nil {
		log.Println(err)
	}

	route, err := mappedTripRouteToProtoRoute(trip_route)

	if err != nil {
		log.Println(err)
	}

	return &pb.PreviewTripResponse{
		Route:     route,
		RideFares: []*pb.RideFare{},
	}, nil
}

func mappedLocationToCoridinate(location *pb.Coordinate) *types.Coordinate {
	return &types.Coordinate{
		Latitude:  location.Latitude,
		Longitude: location.Longtitude,
	}
}

func mappedTripRouteToProtoRoute(trip_route *types.OsrmApiResponse) (*pb.Route, error) {
	var route pb.Route

	if err := mergo.Merge(&route, trip_route.Routes); err != nil {
		return nil, err
	}
	return &route, nil
}
