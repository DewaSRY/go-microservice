package domain

import (
	"context"
	"ride-sharing/services/trip-service/pkg/types"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *types.TripModel) (*types.TripModel, error)
	SaveRideFare(ctx context.Context, fare *types.RideFareModel) error
	GetFareById(ctx context.Context, fareId string) (*types.RideFareModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *types.RideFareModel) (*types.TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(route *types.OsrmApiResponse) []*types.RideFareModel
	GenerateTripFares(ctx context.Context, fares []*types.RideFareModel, userId string, route *types.OsrmApiResponse) ([]*types.RideFareModel, error)
	GetFareById(ctx context.Context, fareId string) (*types.RideFareModel, error)
}
