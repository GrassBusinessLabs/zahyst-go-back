package requests

import (
	"boilerplate/internal/domain"
)

type CreateLocationRequest struct {
	Type        string  `json:"type" validate:"required"`
	Address     string  `json:"address" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Lat         float64 `json:"lat" validate:"required"`
	Lon         float64 `json:"lon" validate:"required"`
}

type UpdateLocationRequest struct {
	Type        string  `json:"type" validate:"required"`
	Address     string  `json:"address" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Lat         float64 `json:"lat" validate:"required"`
	Lon         float64 `json:"lon" validate:"required"`
}

type FindByAreaLocationRequest struct {
	UpperLeftPoint   map[string]float32 `json:"upper_left_point" validate:"required"`
	BottomRightPoint map[string]float32 `json:"bottom_right_point" validate:"required"`
}

func (r CreateLocationRequest) ToDomainModel() (interface{}, error) {
	return domain.Location{
		Type:        r.Type,
		Address:     r.Address,
		Title:       r.Title,
		Description: r.Description,
		Lat:         r.Lat,
		Lon:         r.Lon,
	}, nil
}

func (r UpdateLocationRequest) ToDomainModel() (interface{}, error) {
	return domain.Location{
		Type:        r.Type,
		Address:     r.Address,
		Title:       r.Title,
		Description: r.Description,
		Lat:         r.Lat,
		Lon:         r.Lon,
	}, nil
}
