package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s Service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		Id:       primitive.NewObjectID(),
		UserId:   fare.UserId,
		Status:   "pending",
		RideFare: *fare,
	}
	return s.repo.CreateTrip(ctx, t)
}

func (s Service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson", pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)

	fmt.Println(url)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed_to_get_the_route: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("failed_to_read_the_response: %v", err)
	}

	var routeResp types.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed_to_parse_response: %v", err)
	}

	return &routeResp, nil

}

func (s *Service) EstimatePackagesPriceWithRoute(route *types.OsrmApiResponse) []*domain.RideFareModel {
	baseFares := getBaseFares()

	fareList := make([]*domain.RideFareModel, len(baseFares))

	for i, f := range baseFares {
		fareList[i] = estimationFareRoute(f, route)
	}

	return baseFares
}

func (s *Service) GenerateTripFares(ctx context.Context, fares []*domain.RideFareModel, userId string) ([]*domain.RideFareModel, error) {
	faresList := make([]*domain.RideFareModel, len(fares))

	for i, f := range fares {
		Id := primitive.NewObjectID()
		fare := &domain.RideFareModel{
			UserId:            userId,
			Id:                Id,
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug:       f.PackageSlug,
		}
		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("failed_to_save_trip_fare :%v", err)
		}
		faresList[i] = fare
	}
	return faresList, nil
}

func estimationFareRoute(f *domain.RideFareModel, route *types.OsrmApiResponse) *domain.RideFareModel {
	pricingCfg := types.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents

	distanceKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	distanceFare := distanceKm * pricingCfg.PricePerUnitOfDistance
	timeFare := durationInMinutes * pricingCfg.PricingPerMinute
	totalPrice := carPackagePrice + distanceFare + timeFare

	return &domain.RideFareModel{
		TotalPriceInCents: totalPrice,
		PackageSlug:       f.PackageSlug,
	}
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
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
