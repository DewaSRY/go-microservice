package types

import (
	tripPb "ride-sharing/shared/proto/trip"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	Id       primitive.ObjectID
	UserId   string
	Status   string
	RideFare RideFareModel
	Driver   *TripDriver
}

func (t *TripModel) ToTripProtoTrip() *tripPb.Trip {

	return &tripPb.Trip{
		Id:           t.Id.Hex(),
		UserID:       t.UserId,
		Status:       t.Status,
		SelectedFare: t.RideFare.ToProtoRideFare(),
		Route:        t.RideFare.Route.ToTripProtoRoute(),
		Driver:       t.Driver.ToTripProtoDriver(),
	}
}
