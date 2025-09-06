package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/pkg/types"
	pbd "ride-sharing/shared/proto/driver"
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

func (r *InMemoryRepository) GetTripByID(ctx context.Context, id string) (*types.TripModel, error) {
	trip, ok := r.trip[id]
	if !ok {
		return nil, nil
	}
	return trip, nil
}

func (r *InMemoryRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error {
	trip, ok := r.trip[tripID]
	if !ok {
		return fmt.Errorf("trip not found with ID: %s", tripID)
	}

	trip.Status = status

	if driver != nil {
		trip.Driver = &types.TripDriver{
			Id:             driver.Id,
			Name:           driver.Name,
			ProfilePicture: driver.ProfilePicture,
			CartPlate:      driver.CarPlate,
		}
	}
	return nil
}
