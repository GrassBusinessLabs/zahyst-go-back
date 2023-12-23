package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"encoding/json"
	"log"
	"net/http"
)

type LocationController struct {
	locationService app.LocationService
}

func NewLocationController(ls app.LocationService) LocationController {
	return LocationController{
		locationService: ls,
	}
}

func (c LocationController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		location, err := requests.Bind(r, requests.CreateLocationRequest{}, domain.Location{})
		if err != nil {
			log.Printf("LocationController: %s", err)
			BadRequest(w, err)
			return
		}
		location.UserId = r.Context().Value(UserKey).(domain.User).Id
		location, err = c.locationService.Save(location)
		if err != nil {
			log.Printf("LocationController: %s", err)
			BadRequest(w, err)
			return
		}
		var locationDto resources.LocationDto
		Created(w, locationDto.DomainToDto(location))
	}
}

func (c LocationController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		location, err := requests.Bind(r, requests.UpdateLocationRequest{}, domain.Location{})
		if err != nil {
			log.Printf("LocationController: %s", err)
			BadRequest(w, err)
			return
		}
		instance := r.Context().Value(LocationKey).(domain.Location)
		location.UserId = instance.UserId
		location.Id = instance.Id
		location, err = c.locationService.Update(location)
		if err != nil {
			log.Printf("LocationController: %s", err)
			InternalServerError(w, err)
			return
		}
		var locationDto resources.LocationDto
		Success(w, locationDto.DomainToDto(location))
	}
}

func (c LocationController) Detail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		location := r.Context().Value(LocationKey).(domain.Location)
		var locationDto resources.LocationDto
		Success(w, locationDto.DomainToDto(location))
	}
}

func (c LocationController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locationId := r.Context().Value(LocationKey).(domain.Location).Id
		err := c.locationService.Delete(locationId)
		if err != nil {
			log.Printf("LocationController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (c LocationController) FindByArea() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("LocationController: %s", err)
			InternalServerError(w, err)
			return
		}
		req := requests.FindByAreaLocationRequest{}
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("LocationController: %s", err)
			return
		}
		points := make(map[string][]map[string]float32)
		points["UpperLeftPoint"] = []map[string]float32{{"lat": req.UpperLeftPoint["lat"]}, {"lon": req.UpperLeftPoint["lon"]}}
		points["BottomRightPoint"] = []map[string]float32{{"lat": req.BottomRightPoint["lat"]}, {"lon": req.BottomRightPoint["lon"]}}
		locations, err := c.locationService.FindByArea(pagination, points)
		if err != nil {
			log.Printf("LocationController: %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, resources.LocationDto{}.DomainToDtoPaginatedCollection(locations, pagination))
	}
}

func (c LocationController) FindByUserId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("LocationController: %s", err)
			InternalServerError(w, err)
			return
		}
		userId := r.Context().Value(UserKey).(domain.User).Id
		locations, err := c.locationService.FindByUserId(pagination, userId)
		if err != nil {
			log.Printf("LocationController: %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, resources.LocationDto{}.DomainToDtoPaginatedCollection(locations, pagination))
	}
}
