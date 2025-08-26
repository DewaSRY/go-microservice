package domain

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
}
