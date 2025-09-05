package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/pkg/types"
)

type InMemoryRepository struct {
	trip      map[string]*types.TripModel
	rideFares map[string]*types.RideFareModel
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		trip:      make(map[string]*types.TripModel),
		rideFares: make(map[string]*types.RideFareModel),
	}
}

func (r *InMemoryRepository) CreateTrip(ctx context.Context, trip *types.TripModel) (*types.TripModel, error) {
	r.trip[trip.Id.Hex()] = trip
	return trip, nil
}

func (r *InMemoryRepository) SaveRideFare(ctx context.Context, fare *types.RideFareModel) error {
	r.rideFares[fare.Id.Hex()] = fare
	return nil
}

func (r *InMemoryRepository) GetFareById(ctx context.Context, fareId string) (*types.RideFareModel, error) {

	fare, exists := r.rideFares[fareId]

	if !exists {
		return nil, fmt.Errorf("fare_not_fount")
	}

	return fare, nil

}
