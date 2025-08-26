package main

import (
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

type PreviewTripRequest struct {
	UserID      string           `json:userId`
	Pickup      types.Coordinate `json:pickup`
	Destination types.Coordinate `json:destination`
}

func (p *PreviewTripRequest) mappedToProto() *pb.PreviewTripRequest {

	return &pb.PreviewTripRequest{
		UserID:        p.UserID,
		EndLocation:   mappedLocationToCoridinate(&p.Destination),
		StartLocation: mappedLocationToCoridinate(&p.Pickup),
	}
}

func mappedLocationToCoridinate(location *types.Coordinate) *pb.Coordinate {
	return &pb.Coordinate{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
}
