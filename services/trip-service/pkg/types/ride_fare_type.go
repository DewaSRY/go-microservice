package types

import (
	"time"

	tripPb "ride-sharing/shared/proto/trip"

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

func (t *RideFareModel) ToProtoRideFare() *tripPb.RideFare {
	return &tripPb.RideFare{
		Id:                t.Id.Hex(),
		UserID:            t.UserId,
		PackageSlug:       t.PackageSlug,
		TotalPriceInCents: t.TotalPriceInCents,
	}
}

func GetBaseFares() []*RideFareModel {
	return []*RideFareModel{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 200.0,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 350.0,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 400.0,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000.0,
		},
	}
}
