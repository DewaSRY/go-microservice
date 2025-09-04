package mapper

import (
	"ride-sharing/shared/types"

	driverPb "ride-sharing/shared/proto/driver"
	pb "ride-sharing/shared/proto/trip"
)

func MappedLocationToCoridinate(location *pb.Coordinate) *types.Coordinate {
	return &types.Coordinate{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
}

func MappedGeometryToProtoGeometry(geometry *types.Geometry) *pb.Geometry {
	cordinates := make([]*pb.Coordinate, len(geometry.Coordinates))

	for i, c := range geometry.Coordinates {
		cordinates[i] = &pb.Coordinate{
			Latitude:  c[0],
			Longitude: c[1],
		}
	}

	return &pb.Geometry{
		Coordinates: cordinates,
	}
}

func MappedRouteToProtoroute(route *types.Routes) *pb.Route {

	return &pb.Route{
		Geometry: MappedGeometryToProtoGeometry(&route.Geometry),
		Distance: route.Distance,
		Duration: route.Duration,
	}
}

func MappedTripModelToProtoTripModel(trip *types.TripModel) *pb.Trip {
	return &pb.Trip{
		Id:           trip.Id.Hex(),
		UserID:       trip.UserId,
		Status:       trip.Status,
		SelectedFare: MappedRideFareToProtoRideFare(&trip.RideFare),
		Route:        MappedRouteToProtoroute(&trip.RideFare.Route),
		Driver:       MappedTripDriverToProtoTripDriver(trip.Driver),
	}
}

func MappedRideFareToProtoRideFare(rideFare *types.RideFareModel) *pb.RideFare {
	return &pb.RideFare{
		Id:               rideFare.Id.Hex(),
		UserID:           rideFare.UserId,
		PackageSlug:      rideFare.PackageSlug,
		TotalPriceInCets: rideFare.TotalPriceInCents,
	}
}

func MappedTripDriverToProtoTripDriver(driver *types.TripDriver) *pb.TripDriver {
	return &pb.TripDriver{
		Id:             driver.Id,
		Name:           driver.Name,
		ProfilePicture: driver.ProfilePicture,
	}
}

func MappedCoridinateToPortoLatlong(coordinate *types.Coordinate) *driverPb.Location {
	return &driverPb.Location{
		Latitude:  coordinate.Latitude,
		Longitude: coordinate.Longitude,
	}
}
