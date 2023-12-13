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
	Lat1 float32 `json:"lat1" validate:"required"`
	Lon1 float32 `json:"lon1" validate:"required"`
	Lat2 float32 `json:"lat2" validate:"required"`
	Lon2 float32 `json:"lon2" validate:"required"`
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
		Lat1: r.Lat1,
		Lon1: r.Lon1,
		Lat2: r.Lat2,
		Lon2: r.Lon2,
	}, nil
}
