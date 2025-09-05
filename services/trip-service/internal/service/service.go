package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"

	"ride-sharing/services/trip-service/pkg/types"

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

func (s Service) CreateTrip(ctx context.Context, fare *types.RideFareModel) (*types.TripModel, error) {
	t := &types.TripModel{
		Id:       primitive.NewObjectID(),
		UserId:   fare.UserId,
		Status:   "pending",
		RideFare: *fare,
		Driver:   &types.TripDriver{},
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

func (s *Service) EstimatePackagesPriceWithRoute(route *types.OsrmApiResponse) []*types.RideFareModel {
	baseFares := types.GetBaseFares()

	fareList := make([]*types.RideFareModel, len(baseFares))

	for i, f := range baseFares {
		fareList[i] = estimationFareRoute(f, route)
	}

	return baseFares
}

func (s *Service) GenerateTripFares(ctx context.Context, fares []*types.RideFareModel, userId string, route *types.OsrmApiResponse) ([]*types.RideFareModel, error) {
	faresList := make([]*types.RideFareModel, len(fares))

	for i, f := range fares {
		Id := primitive.NewObjectID()
		fare := &types.RideFareModel{
			UserId:            userId,
			Id:                Id,
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug:       f.PackageSlug,
			Route:             route.Routes[0],
		}

		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("failed_to_save_trip_fare :%v", err)
		}

		faresList[i] = fare
	}
	return faresList, nil
}

func estimationFareRoute(f *types.RideFareModel, route *types.OsrmApiResponse) *types.RideFareModel {
	pricingCfg := types.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents

	distanceKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	distanceFare := distanceKm * pricingCfg.PricePerUnitOfDistance
	timeFare := durationInMinutes * pricingCfg.PricingPerMinute
	totalPrice := carPackagePrice + distanceFare + timeFare

	return &types.RideFareModel{
		TotalPriceInCents: totalPrice,
		PackageSlug:       f.PackageSlug,
	}
}

func (r *Service) GetFareById(ctx context.Context, fareId string) (*types.RideFareModel, error) {

	return r.repo.GetFareById(ctx, fareId)
}
