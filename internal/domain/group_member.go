package domain

import "time"

type GroupMember struct {
	Id          uint64
	UserId      uint64
	GroupId     uint64
	AccessLevel string
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type GroupMembers struct {
	Items []GroupMember
	Total uint64
	Pages uint
}

func (groupMember GroupMember) GetUserId() uint64 {
	return groupMember.UserId
}
