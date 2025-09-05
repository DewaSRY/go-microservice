package types

// type Route struct {
// 	Distance float64     `json:"distance"`
// 	Duration float64     `json:"duration"`
// 	Geometry []*Geometry `json:"geometry"`
// }

// type Coordinate struct {
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// }

// type Coordinates [2]float64

// type OsrmApiResponse struct {
// 	Code      string      `json:"code"`
// 	Routes    []Routes    `json:"routes"`
// 	Waypoints []Waypoints `json:"waypoints"`
// }

// type Legs struct {
// 	Steps    []any   `json:"steps"`
// 	Weight   float64 `json:"weight"`
// 	Summary  string  `json:"summary"`
// 	Duration float64 `json:"duration"`
// 	Distance float64 `json:"distance"`
// }

// type Geometry struct {
// 	Coordinates []Coordinates `json:"coordinates"`
// 	Type        string        `json:"type"`
// }

// type Routes struct {
// 	Legs       []Legs   `json:"legs"`
// 	WeightName string   `json:"weight_name"`
// 	Geometry   Geometry `json:"geometry"`
// 	Weight     float64  `json:"weight"`
// 	Duration   float64  `json:"duration"`
// 	Distance   float64  `json:"distance"`
// }

// type Waypoints struct {
// 	Hint     string    `json:"hint"`
// 	Location []float64 `json:"location"`
// 	Name     string    `json:"name"`
// 	Distance float64   `json:"distance"`
// }

// type PricingConfig struct {
// 	PricePerUnitOfDistance float64
// 	PricingPerMinute       float64
// }

// func DefaultPricingConfig() *PricingConfig {
// 	return &PricingConfig{
// 		PricePerUnitOfDistance: 1.5,
// 		PricingPerMinute:       0.25,
// 	}
// }
