package database

import (
	"boilerplate/internal/domain"
	"errors"
	"math"
	"time"

	"github.com/upper/db/v4"
)

const LocationsTableName = "locations"

type location struct {
	Id          uint64     `db:"id,omitempty"`
	UserId      uint64     `db:"user_id,omitempty"`
	Type        string     `db:"type"`
	Address     string     `db:"address"`
	Lat         float64    `db:"lat"`
	Lon         float64    `db:"lon"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type LocationRepository interface {
	Save(sess domain.Location) (domain.Location, error)
	Update(location domain.Location, userId uint64) (domain.Location, error)
	Delete(id uint64, userId uint64) error
	Detail(id uint64) (domain.Location, error)
	FindByArea(p domain.Pagination, area_points domain.AreaPoints) (domain.Locations, error)
	FindByUserId(p domain.Pagination, user_id uint64) (domain.Locations, error)
	FindById(id uint64) (domain.Location, error)
}

type locationRepository struct {
	coll db.Collection
}

func NewLocationRepository(dbSession db.Session) locationRepository {
	return locationRepository{
		coll: dbSession.Collection(LocationsTableName),
	}
}

func (r locationRepository) Save(location domain.Location) (domain.Location, error) {
	loc := r.mapDomainToModel(location)
	loc.CreatedDate, loc.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&loc)
	if err != nil {
		return domain.Location{}, err
	}
	return r.mapModelToDomain(loc), nil
}

func (r locationRepository) Update(location domain.Location, userId uint64) (domain.Location, error) {
	loc := r.mapDomainToModel(location)
	loc.UpdatedDate = time.Now()
	if loc.UserId != userId {
		return location, errors.New("Invalid user")
	}
	err := r.coll.Find(db.Cond{"id": loc.Id}).Update(&loc)
	if err != nil {
		return domain.Location{}, err
	}
	return r.mapModelToDomain(loc), nil
}

func (r locationRepository) Detail(id uint64) (domain.Location, error) {
	location, err := r.FindById(id)
	if err != nil {
		return domain.Location{}, err
	}
	return location, nil
}

func (r locationRepository) Delete(id uint64, userId uint64) error {
	location := r.coll.Find(db.Cond{"id": id, "user_id": userId, "deleted_date": nil})
	exists, err := location.Exists()
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("No such location")
	}
	return location.Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r locationRepository) FindByArea(p domain.Pagination, area_points domain.AreaPoints) (domain.Locations, error) {
	var data []location
	query := r.coll.Find(db.Cond{"lat >": area_points.Lat1, "lat <": area_points.Lat2, "lon <": area_points.Lon1, "lon >": area_points.Lon2})
	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.Locations{}, err
	}

	locations := r.mapModelToDomainPagination(data)

	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Locations{}, err
	}

	locations.Total = totalCount
	locations.Pages = uint(math.Ceil(float64(locations.Total) / float64(p.CountPerPage)))

	return locations, nil
}

func (r locationRepository) FindByUserId(p domain.Pagination, user_id uint64) (domain.Locations, error) {
	var data []location
	query := r.coll.Find(db.Cond{"user_id": user_id})
	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.Locations{}, err
	}

	locations := r.mapModelToDomainPagination(data)

	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Locations{}, err
	}

	locations.Total = totalCount
	locations.Pages = uint(math.Ceil(float64(locations.Total) / float64(p.CountPerPage)))

	return locations, nil
}

func (r locationRepository) FindById(id uint64) (domain.Location, error) {
	var loc location
	err := r.coll.Find(db.Cond{"id": id}).One(&loc)
	if err != nil {
		return domain.Location{}, err
	}
	return r.mapModelToDomain(loc), nil
}

func (r locationRepository) mapDomainToModel(d domain.Location) location {
	return location{
		Id:          d.Id,
		UserId:      d.UserId,
		Type:        d.Type,
		Address:     d.Address,
		Lat:         d.Lat,
		Lon:         d.Lon,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r locationRepository) mapModelToDomain(m location) domain.Location {
	return domain.Location{
		Id:          m.Id,
		UserId:      m.UserId,
		Type:        m.Type,
		Address:     m.Address,
		Lat:         m.Lat,
		Lon:         m.Lon,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (f locationRepository) mapModelToDomainPagination(locations []location) domain.Locations {
	new_locations := make([]domain.Location, len(locations))
	for i, location := range locations {
		new_locations[i] = f.mapModelToDomain(location)
	}
	return domain.Locations{Items: new_locations}
}
