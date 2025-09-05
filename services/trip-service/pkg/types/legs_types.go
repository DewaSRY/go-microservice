package types

type Legs struct {
	Steps    []any   `json:"steps"`
	Weight   float64 `json:"weight"`
	Summary  string  `json:"summary"`
	Duration float64 `json:"duration"`
	Distance float64 `json:"distance"`
}
