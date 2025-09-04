package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/events"
	pb "ride-sharing/shared/proto/trip"

	"ride-sharing/shared/mapper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service   domain.TripService
	publisher *events.TripEventPublisher
}

func NewGRPCHandler(
	server *grpc.Server,
	service domain.TripService,
	publisher *events.TripEventPublisher) *GRPCHandler {

	handler := &GRPCHandler{
		service:   service,
		publisher: publisher,
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
		return nil, status.Errorf(codes.Internal, "failed_to_create_route")
	}

	userId := req.GetUserID()

	estimatedFares := h.service.EstimatePackagesPriceWithRoute(trip_route)

	fares, err := h.service.GenerateTripFares(ctx, estimatedFares, userId, trip_route)

	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed_to_generate_trip_fare")
	}

	rideFare := make([]*pb.RideFare, len(fares))
	for i, f := range fares {
		rideFare[i] = mapper.MappedRideFareToProtoRideFare(f)
	}

	return &pb.PreviewTripResponse{
		Route:     mapper.MappedRouteToProtoroute(&trip_route.Routes[0]),
		RideFares: rideFare,
	}, nil
}

func (h *GRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fareID := req.GetRideFareID()
	userID := req.GetUserID()
	// 1. Fetch and validate the fare.
	fare, err := h.service.GetFareById(ctx, fareID)

	if err != nil || fare.UserId != userID {
		return nil, status.Errorf(codes.NotFound, "fare_not_found")
	}

	// 2. Call create trip
	trip, err := h.service.CreateTrip(ctx, fare)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed_to_create_trip")
	}
	// 3. We also need to initialize an empty drver to the trip.
	// 4. Add a comment at the end of the function to publish an event on the Asnyc Comms module.
	if err := h.publisher.PublishTripCreated(ctx, *trip); err != nil {
		return nil, status.Errorf(codes.Internal, "failed_to_publsh_the_trip_created_event: %v", err)
	}

	return &pb.CreateTripResponse{
		TripID: trip.Id.Hex(),
	}, nil
}
