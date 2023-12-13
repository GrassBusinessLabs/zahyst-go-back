package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"log"
)

type LocationService interface {
	Save(location domain.Location) (domain.Location, error)
	Update(location domain.Location, userId uint64) (domain.Location, error)
	Delete(id uint64, userId uint64) error
	Detail(id uint64) (domain.Location, error)
	FindByArea(p domain.Pagination, area_points domain.AreaPoints) (domain.Locations, error)
	FindByUserId(p domain.Pagination, user_id uint64) (domain.Locations, error)
	FindById(id uint64) (domain.Location, error)
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

func (s locationService) Update(location domain.Location, userId uint64) (domain.Location, error) {
	loc, err := s.locationRepo.Update(location, userId)
	if err != nil {
		log.Printf("LocationService: %s", err)
		return loc, err
	}

	return loc, err
}

func (s locationService) Detail(id uint64) (domain.Location, error) {
	location, err := s.locationRepo.Detail(id)
	if err != nil {
		return domain.Location{}, err
	}

	return location, nil
}

func (s locationService) Delete(id uint64, userId uint64) error {
	err := s.locationRepo.Delete(id, userId)
	if err != nil {
		log.Printf("LocationService: %s", err)
		return err
	}

	return nil
}

func (s locationService) FindByArea(p domain.Pagination, area_points domain.AreaPoints) (domain.Locations, error) {
	locations, err := s.locationRepo.FindByArea(p, area_points)
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

func (s locationService) FindById(id uint64) (domain.Location, error) {
	location, err := s.locationRepo.FindById(id)
	if err != nil {
		log.Printf("LocationService: %s", err)
		return domain.Location{}, err
	}

	return location, err
}
