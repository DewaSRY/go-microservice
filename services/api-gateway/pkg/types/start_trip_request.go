package types

import (
	tripPb "ride-sharing/shared/proto/trip"
)

type StartTripRequest struct {
	RideFareID string `json:"rideFareID"`
	UserID     string `json:"userID"`
}

func (t *StartTripRequest) ToTripProtoCreateTrip() *tripPb.CreateTripRequest {
	return &tripPb.CreateTripRequest{
		RideFareID: t.RideFareID,
		UserID:     t.UserID,
	}
}
