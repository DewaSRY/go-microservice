package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	Id                primitive.ObjectID
	UserId            string
	PackageSlug       string
	TotalPriceInCents float64
	Expires           time.Time
	Route             Routes
}

type TripModel struct {
	Id       primitive.ObjectID
	UserId   string
	Status   string
	RideFare RideFareModel
	Driver   *TripDriver
}

type TripDriver struct {
	Id             string
	Name           string
	ProfilePicture string
	CartPlate      string
}
