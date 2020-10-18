package gnis

// Feature .
type Feature struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Class     string  `json:"class"`
	State     string  `json:"state"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation int     `json:"elevation"`
}
