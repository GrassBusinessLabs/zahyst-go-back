package resources

import (
	"boilerplate/internal/domain"
)

type LocationDto struct {
	Id          uint64  `json:"id,omitempty"`
	UserId      uint64  `json:"user_id"`
	Type        string  `json:"type"`
	Address     string  `json:"address"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

type LocationsDto struct {
	Items []LocationDto `json:"items"`
	Total uint64        `json:"total"`
	Pages uint          `json:"pages"`
}

func (d LocationDto) DomainToDto(location domain.Location) LocationDto {
	return LocationDto{
		Id:          location.Id,
		UserId:      location.UserId,
		Type:        location.Type,
		Address:     location.Address,
		Title:       location.Title,
		Description: location.Description,
		Lat:         location.Lat,
		Lon:         location.Lon,
	}
}

func (d LocationDto) DomainToDtoCollection(locations domain.Locations) LocationsDto {
	result := make([]LocationDto, len(locations.Items))

	for i := range locations.Items {
		result[i] = d.DomainToDto(locations.Items[i])
	}

	return LocationsDto{Items: result, Pages: locations.Pages, Total: locations.Total}
}

func (d LocationDto) DomainToDtoPaginatedCollection(locations domain.Locations, pag domain.Pagination) LocationsDto {
	result := make([]LocationDto, len(locations.Items))

	for i := range locations.Items {
		result[i] = d.DomainToDto(locations.Items[i])
	}

	return LocationsDto{Items: result, Pages: locations.Pages, Total: locations.Total}
}
