package types

import (
	driverPb "ride-sharing/shared/proto/driver"
	tripPb "ride-sharing/shared/proto/trip"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewCoordinateFormProto(location *tripPb.Coordinate) *Coordinate {
	return &Coordinate{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
}

func (t *Coordinate) ToDriverProtoLocation() *driverPb.Location {
	return &driverPb.Location{
		Latitude:  t.Latitude,
		Longitude: t.Longitude,
	}
}
