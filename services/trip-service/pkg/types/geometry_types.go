package types

import (
	tripPb "ride-sharing/shared/proto/trip"
)

type Geometry struct {
	Coordinates []Coordinates `json:"coordinates"`
	Type        string        `json:"type"`
}

func (t *Geometry) ToTripProtoGeometry() *tripPb.Geometry {
	coordinates := make([]*tripPb.Coordinate, len(t.Coordinates))

	for i, c := range t.Coordinates {
		coordinates[i] = &tripPb.Coordinate{
			Latitude:  c[0],
			Longitude: c[1],
		}
	}

	return &tripPb.Geometry{
		Coordinates: coordinates,
	}
}
