package types

import (
	tripPb "ride-sharing/shared/proto/trip"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (t *Coordinate) ToTripProtoCoordinate() *tripPb.Coordinate {
	return &tripPb.Coordinate{
		Latitude:  t.Latitude,
		Longitude: t.Longitude,
	}
}
