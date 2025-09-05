package types

import (
	tripPb "ride-sharing/shared/proto/trip"
)

type Routes struct {
	Legs       []Legs   `json:"legs"`
	WeightName string   `json:"weight_name"`
	Geometry   Geometry `json:"geometry"`
	Weight     float64  `json:"weight"`
	Duration   float64  `json:"duration"`
	Distance   float64  `json:"distance"`
}

func (t *Routes) ToTripProtoRoute() *tripPb.Route {
	return &tripPb.Route{
		Geometry: t.Geometry.ToTripProtoGeometry(),
		Distance: t.Distance,
		Duration: t.Duration,
	}
}
