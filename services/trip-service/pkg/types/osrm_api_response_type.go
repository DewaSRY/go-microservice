package types

type OsrmApiResponse struct {
	Code      string      `json:"code"`
	Routes    []Routes    `json:"routes"`
	Waypoints []Waypoints `json:"waypoints"`
}
