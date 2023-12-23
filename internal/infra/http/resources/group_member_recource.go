package resources

import (
	"boilerplate/internal/domain"
)

type GroupMemberDto struct {
	Id          uint64 `json:"id,omitempty"`
	UserId      uint64 `json:"user_id"`
	GroupId     uint64 `json:"group_id"`
	AccessLevel string `json:"access_level"`
}

type GroupMembersDto struct {
	Items []GroupMemberDto `json:"items"`
	Total uint64           `json:"total"`
	Pages uint             `json:"pages"`
}

func (d GroupMemberDto) DomainToDto(groupMember domain.GroupMember) GroupMemberDto {
	return GroupMemberDto{
		Id:          groupMember.Id,
		UserId:      groupMember.UserId,
		GroupId:     groupMember.GroupId,
		AccessLevel: groupMember.AccessLevel,
	}
}

func (d GroupMemberDto) DomainToDtoCollection(groupMembers domain.GroupMembers) GroupMembersDto {
	result := make([]GroupMemberDto, len(groupMembers.Items))

	for i := range groupMembers.Items {
		result[i] = d.DomainToDto(groupMembers.Items[i])
	}

	return GroupMembersDto{Items: result, Pages: groupMembers.Pages, Total: groupMembers.Total}
}

func (d GroupMemberDto) DomainToDtoPaginatedCollection(groupMembers domain.GroupMembers, pag domain.Pagination) GroupMembersDto {
	result := make([]GroupMemberDto, len(groupMembers.Items))

	for i := range groupMembers.Items {
		result[i] = d.DomainToDto(groupMembers.Items[i])
	}

	return GroupMembersDto{Items: result, Pages: groupMembers.Pages, Total: groupMembers.Total}
}
