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

	userId := req.GetUserID()

	routes := make([]*pb.Route, len(trip_route.Routes))

	for i, r := range trip_route.Routes {
		routes[i] = mapper.MappedRouteToProtoroute(&r)
	}

	if err != nil {
		log.Println(err)
	}

	estimatedFares := h.service.EstimatePackagesPriceWithRoute(trip_route)

	fares, err := h.service.GenerateTripFares(ctx, estimatedFares, userId)

	if err != nil {
		log.Println(err)
	}

	rideFare := make([]*pb.RideFare, len(fares))
	for i, f := range fares {
		rideFare[i] = mappedRideFareToProtoRideFare(f)
	}

	return &pb.PreviewTripResponse{
		Route:     routes,
		RideFares: rideFare,
	}, nil
}

func mappedRideFareToProtoRideFare(fare *domain.RideFareModel) *pb.RideFare {
	return &pb.RideFare{
		Id:               fare.Id.Hex(),
		UserID:           fare.UserId,
		PackageSlug:      fare.PackageSlug,
		TotalPriceInCets: fare.TotalPriceInCents,
	}
}
