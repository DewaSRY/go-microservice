package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"
)

type InMemoryRepository struct {
	trip      map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		trip:      make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *InMemoryRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trip[trip.Id.Hex()] = trip
	return trip, nil
}

func (s *InMemoryRepository) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson", pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed_to_get_the_route: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("failed_to_read_the_response: %v", err)
	}

	var routeResp types.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed_to_parse_response: %v", err)
	}

	return &routeResp, nil

}

func (r *InMemoryRepository) SaveRideFare(ctx context.Context, fare *domain.RideFareModel) error {
	r.rideFares[fare.Id.Hex()] = fare
	return nil
}

func (r *InMemoryRepository) GetFareById(ctx context.Context, fareId string) (*domain.RideFareModel, error) {

	fare, exists := r.rideFares[fareId]

	if !exists {
		return nil, fmt.Errorf("fare_not_fount")
	}

	return fare, nil

}
