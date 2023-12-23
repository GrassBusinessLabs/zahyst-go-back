package domain

import "time"

type Group struct {
	Id          uint64
	Title       string
	Description string
	UserId      uint64
	AcessCode   string
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type Groups struct {
	Items []Group
	Total uint64
	Pages uint
}

func (group Group) GetUserId() uint64 {
	return group.UserId
}
