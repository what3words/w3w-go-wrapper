package core

import "fmt"

// Coordinates models a starndard struct
// to represent a valid coordinate pair
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (c Coordinates) String() string {
	return fmt.Sprintf("%.6f,%.6f", c.Lat, c.Lng)
}
