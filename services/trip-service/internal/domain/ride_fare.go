package domain

import (
	"context"
	"ride-sharing/shared/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	Id                primitive.ObjectID
	UserId            string
	PackageSlug       string
	totalPriceInCents float64
	Expires           time.Time
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmApiResponse, error)
}
