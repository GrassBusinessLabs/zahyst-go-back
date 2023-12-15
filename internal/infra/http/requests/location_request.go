package requests

import (
	"boilerplate/internal/domain"
)

type CreateLocationRequest struct {
	Type    string  `json:"type" validate:"required"`
	Address string  `json:"address" validate:"required"`
	Lat     float64 `json:"lat" validate:"required"`
	Lon     float64 `json:"lon" validate:"required"`
}

type UpdateLocationRequest struct {
	Type    string  `json:"type" validate:"required"`
	Address string  `json:"address" validate:"required"`
	Lat     float64 `json:"lat" validate:"required"`
	Lon     float64 `json:"lon" validate:"required"`
}

type FindByAreaLocationRequest struct {
	UpperLeftPoint   map[string]float32 `json:"upper_left_point" validate:"required"`
	BottomRightPoint map[string]float32 `json:"bottom_right_point" validate:"required"`
}

func (r CreateLocationRequest) ToDomainModel() (interface{}, error) {
	return domain.Location{
		Type:    r.Type,
		Address: r.Address,
		Lat:     r.Lat,
		Lon:     r.Lon,
	}, nil
}

func (r UpdateLocationRequest) ToDomainModel() (interface{}, error) {
	return domain.Location{
		Type:    r.Type,
		Address: r.Address,
		Lat:     r.Lat,
		Lon:     r.Lon,
	}, nil
}

func (r FindByAreaLocationRequest) ToDomainModel() (interface{}, error) {
	return domain.AreaPoints{
		Lat1: r.UpperLeftPoint["lat"],
		Lon1: r.UpperLeftPoint["lon"],
		Lat2: r.BottomRightPoint["lat"],
		Lon2: r.BottomRightPoint["lon"],
	}, nil
}
