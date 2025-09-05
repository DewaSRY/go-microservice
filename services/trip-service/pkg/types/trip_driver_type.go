package types

import (
	tripPb "ride-sharing/shared/proto/trip"
)

type TripDriver struct {
	Id             string
	Name           string
	ProfilePicture string
	CartPlate      string
}

func (t *TripDriver) ToTripProtoDriver() *tripPb.TripDriver {
	return &tripPb.TripDriver{
		Id:             t.Id,
		Name:           t.Name,
		ProfilePicture: t.ProfilePicture,
	}
}
