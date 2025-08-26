package mapper

import (
	"ride-sharing/shared/types"

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
