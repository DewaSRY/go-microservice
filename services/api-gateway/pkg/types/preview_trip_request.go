package types

import (
	tripPb "ride-sharing/shared/proto/trip"
)

type PreviewTripRequest struct {
	UserID      string     `json:userId`
	Pickup      Coordinate `json:pickup`
	Destination Coordinate `json:destination`
}

func (t *PreviewTripRequest) ToTripProtoTripRequest() *tripPb.PreviewTripRequest {
	return &tripPb.PreviewTripRequest{
		UserID:        t.UserID,
		EndLocation:   t.Destination.ToTripProtoCoordinate(),
		StartLocation: t.Pickup.ToTripProtoCoordinate(),
	}
}
