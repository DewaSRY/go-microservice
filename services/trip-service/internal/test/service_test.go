package test

import (
	"context"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/mapper"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
	"testing"
)

func TestService(t *testing.T) {
	ctx, cencel := context.WithCancel(context.Background())
	defer cencel()
	repository := repository.NewInMemoryRepository()
	service := service.NewService(repository)

	pickup := &types.Coordinate{
		Latitude:  -8.761044,
		Longitude: 115.158737,
	}

	destination := &types.Coordinate{
		Latitude:  -8.810502,
		Longitude: 115.168345,
	}

	res, err := service.GetRoute(ctx, pickup, destination)

	if err != nil {
		t.Fatalf("failed to get the route: %v", err)
	}

	if err != nil {
		t.Fatalf("failed to get the route: %v", err)
	}

	if res == nil {
		t.Fatalf("expected a route response, got nil")
	}

	routeList := make([]*pb.Route, len(res.Routes))

	for i, r := range res.Routes {
		routeList[i] = mapper.MappedRouteToProtoroute(&r)
	}

	if err != nil {
		t.Fatalf("failed to get the route: %v", err)
	}

	// for _, r := range routeList {
	// 	fmt.Println(r)
	// 	fmt.Println("\n  ")

	// }

}
