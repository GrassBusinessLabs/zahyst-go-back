package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"log"
)

type LocationService interface {
	Save(location domain.Location) (domain.Location, error)
	Update(location domain.Location) (domain.Location, error)
	Delete(id uint64) error
	FindByArea(p domain.Pagination, points map[string][]map[string]float32) (domain.Locations, error)
	FindByUserId(p domain.Pagination, user_id uint64) (domain.Locations, error)
	Find(uint64) (interface{}, error)
}

type locationService struct {
	locationRepo database.LocationRepository
}

func NewLocationService(lr database.LocationRepository) locationService {
	return locationService{
		locationRepo: lr,
	}
}

func (s locationService) Save(location domain.Location) (domain.Location, error) {
	loc, err := s.locationRepo.Save(location)
	if err != nil {
		log.Printf("LocationService: %s", err)
		return domain.Location{}, err
	}

	return loc, err
}

func (s locationService) Update(location domain.Location) (domain.Location, error) {
	loc, err := s.locationRepo.Update(location)
	if err != nil {
		log.Printf("LocationService: %s", err)
		return loc, err
	}

	return loc, err
}

func (s locationService) Delete(id uint64) error {
	err := s.locationRepo.Delete(id)
	if err != nil {
		log.Printf("LocationService: %s", err)
		return err
	}

	return nil
}

func (s locationService) FindByArea(p domain.Pagination, points map[string][]map[string]float32) (domain.Locations, error) {
	locations, err := s.locationRepo.FindByArea(p, points)
	if err != nil {
		log.Printf("LocationService: %s", err)
		return domain.Locations{}, err
	}

	return locations, err
}

func (s locationService) FindByUserId(p domain.Pagination, user_id uint64) (domain.Locations, error) {
	locations, err := s.locationRepo.FindByUserId(p, user_id)
	if err != nil {
		log.Printf("LocationService: %s", err)
		return domain.Locations{}, err
	}

	return locations, err
}

func (s locationService) Find(id uint64) (interface{}, error) {
	location, err := s.locationRepo.FindById(id)
	if err != nil {
		log.Printf("LocationService: %s", err)
		return domain.Location{}, err
	}

	return location, err
}
