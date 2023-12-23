package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"log"
	"net/http"
)

type GroupController struct {
	groupService app.GroupService
}

func NewGroupController(gs app.GroupService) GroupController {
	return GroupController{
		groupService: gs,
	}
}

func (c GroupController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group, err := requests.Bind(r, requests.CreateGroupRequest{}, domain.Group{})
		if err != nil {
			log.Printf("GroupController: %s", err)
			BadRequest(w, err)
			return
		}
		group.UserId = r.Context().Value(UserKey).(domain.User).Id
		group, err = c.groupService.Save(group)
		if err != nil {
			log.Printf("GroupController: %s", err)
			BadRequest(w, err)
			return
		}
		var groupDto resources.GroupDto
		Created(w, groupDto.DomainToDto(group))
	}
}

func (c GroupController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group, err := requests.Bind(r, requests.UpdateGroupRequest{}, domain.Group{})
		if err != nil {
			log.Printf("GroupController: %s", err)
			BadRequest(w, err)
			return
		}
		instance := r.Context().Value(GroupKey).(domain.Group)
		group.UserId = instance.UserId
		group.Id = instance.Id
		group, err = c.groupService.Update(group)
		if err != nil {
			log.Printf("GroupController: %s", err)
			InternalServerError(w, err)
			return
		}
		var groupDto resources.GroupDto
		Success(w, groupDto.DomainToDto(group))
	}
}

func (c GroupController) Detail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.Context().Value(GroupKey).(domain.Group)
		var groupDto resources.GroupDto
		Success(w, groupDto.DomainToDto(group))
	}
}

func (c GroupController) GetList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("GroupController: %s", err)
			InternalServerError(w, err)
			return
		}
		groups, err := c.groupService.GetList(pagination)
		if err != nil {
			log.Printf("GroupController: %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, resources.GroupDto{}.DomainToDtoPaginatedCollection(groups, pagination))
	}
}

func (c GroupController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupId := r.Context().Value(GroupKey).(domain.Group).Id
		err := c.groupService.Delete(groupId)
		if err != nil {
			log.Printf("GroupController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (c GroupController) GetAccessCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.Context().Value(GroupKey).(domain.Group)
		accessCode := c.groupService.GetAccessCode(group)

		Success(w, map[string]string{"accessCode": accessCode})
	}
}
