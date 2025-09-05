package types

type Waypoints struct {
	Hint     string    `json:"hint"`
	Location []float64 `json:"location"`
	Name     string    `json:"name"`
	Distance float64   `json:"distance"`
}
