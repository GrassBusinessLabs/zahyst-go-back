package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GroupMemberController struct {
	groupMemberService app.GroupMemberService
}

func NewGroupMemberController(gms app.GroupMemberService) GroupMemberController {
	return GroupMemberController{
		groupMemberService: gms,
	}
}

func (c GroupMemberController) AddGroupMember() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := requests.AddGroupMemberRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			return
		}
		accessCode := req.AccessCode
		userId := r.Context().Value(UserKey).(domain.User).Id
		groupMember, err := c.groupMemberService.AddGroupMember(accessCode, userId)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			BadRequest(w, err)
			return
		}
		var groupMemberDto resources.GroupMemberDto
		Created(w, groupMemberDto.DomainToDto(groupMember))
	}
}

func (c GroupMemberController) ChangeAccessLevel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupMember := r.Context().Value(GroupMemberKey).(domain.GroupMember)
		req := requests.ChangeMemberAccessLevelRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			return
		}
		newAccessLevel := req.AccessLevel
		groupMember, err = c.groupMemberService.ChangeAccessLevel(groupMember, newAccessLevel)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			InternalServerError(w, err)
			return
		}
		var groupMemberDto resources.GroupMemberDto
		Success(w, groupMemberDto.DomainToDto(groupMember))
	}
}

func (c GroupMemberController) GetMembersList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			InternalServerError(w, err)
			return
		}
		groupId, err := strconv.ParseUint(chi.URLParam(r, "groupId"), 10, 64)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			InternalServerError(w, err)
			return
		}
		groupMembers, err := c.groupMemberService.GetMembersList(pagination, groupId)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, resources.GroupMemberDto{}.DomainToDtoPaginatedCollection(groupMembers, pagination))
	}
}

func (c GroupMemberController) DeleteGroupMember() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupMemberId := r.Context().Value(GroupMemberKey).(domain.GroupMember).Id
		err := c.groupMemberService.DeleteGroupMember(groupMemberId)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (c GroupMemberController) FindMembersByArea() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			InternalServerError(w, err)
			return
		}
		groupId, err := strconv.ParseUint(chi.URLParam(r, "groupId"), 10, 64)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			InternalServerError(w, err)
			return
		}
		req := requests.FindMembersByAreaRequest{}
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			return
		}
		points := make(map[string]map[string]float32)
		points["UpperLeftPoint"] = map[string]float32{"lat": req.UpperLeftPoint["lat"], "lon": req.UpperLeftPoint["lon"]}
		points["BottomRightPoint"] = map[string]float32{"lat": req.BottomRightPoint["lat"], "lon": req.BottomRightPoint["lon"]}
		groupMembers, err := c.groupMemberService.FindMembersByArea(pagination, groupId, points)
		if err != nil {
			log.Printf("GroupMemberController: %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, resources.GroupMemberDto{}.DomainToDtoPaginatedCollection(groupMembers, pagination))
	}
}
