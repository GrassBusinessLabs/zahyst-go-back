package resources

import (
	"boilerplate/internal/domain"
)

type GroupDto struct {
	Id          uint64 `json:"id,omitempty"`
	UserId      uint64 `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GroupsDto struct {
	Items []GroupDto `json:"items"`
	Total uint64     `json:"total"`
	Pages uint       `json:"pages"`
}

func (d GroupDto) DomainToDto(group domain.Group) GroupDto {
	return GroupDto{
		Id:          group.Id,
		UserId:      group.UserId,
		Title:       group.Title,
		Description: group.Description,
	}
}

func (d GroupDto) DomainToDtoCollection(groups domain.Groups) GroupsDto {
	result := make([]GroupDto, len(groups.Items))

	for i := range groups.Items {
		result[i] = d.DomainToDto(groups.Items[i])
	}

	return GroupsDto{Items: result, Pages: groups.Pages, Total: groups.Total}
}

func (d GroupDto) DomainToDtoPaginatedCollection(groups domain.Groups, pag domain.Pagination) GroupsDto {
	result := make([]GroupDto, len(groups.Items))

	for i := range groups.Items {
		result[i] = d.DomainToDto(groups.Items[i])
	}

	return GroupsDto{Items: result, Pages: groups.Pages, Total: groups.Total}
}
