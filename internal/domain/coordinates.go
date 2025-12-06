package domain

import "errors"

var (
	ErrInvalidLatitude  = errors.New("invalid latitude")
	ErrInvalidLongitude = errors.New("invalid longitude")
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewCoordinates(lat, lng float64) (*Coordinates, error) {
	if lat < -90 || lat > 90 {
		return nil, ErrInvalidLatitude
	}
	if lng < -180 || lng > 180 {
		return nil, ErrInvalidLongitude
	}

	return &Coordinates{
		Latitude:  lat,
		Longitude: lng,
	}, nil
}
