package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
			return
		}
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idParam, 0, 64)
		userId := r.Context().Value(UserKey).(domain.User).Id
		location_instance, err := c.locationService.FindById(id)
		location.UserId = location_instance.UserId
		location.Id = id
		location, err = c.locationService.Update(location, userId)
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
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idParam, 0, 64)
		location, err := c.locationService.Detail(id)
		if err != nil {
			log.Printf("LocationController: %s", err)
			InternalServerError(w, err)
			return
		}
		var locationDto resources.LocationDto
		Success(w, locationDto.DomainToDto(location))
	}
}

func (c LocationController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idParam, 0, 64)
		userId := r.Context().Value(UserKey).(domain.User).Id
		err = c.locationService.Delete(id, userId)
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
		area_points, err := requests.Bind(r, requests.FindByAreaLocationRequest{}, domain.AreaPoints{})
		if err != nil {
			log.Printf("LocationController: %s", err)
			InternalServerError(w, err)
			return
		}
		locations, err := c.locationService.FindByArea(pagination, area_points)
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
		userIdParam := chi.URLParam(r, "id")
		userId, err := strconv.ParseUint(userIdParam, 0, 64)
		if err != nil {
			log.Printf("LocationController: %s", err)
			InternalServerError(w, err)
			return
		}
		locations, err := c.locationService.FindByUserId(pagination, userId)
		if err != nil {
			log.Printf("LocationController: %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, resources.LocationDto{}.DomainToDtoPaginatedCollection(locations, pagination))
	}
}
