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

func (groupMember GroupMember) GetAccessLevels() []AccessLevel {
	return []AccessLevel{CasualAccessLevel{}, ModeratorAccessLevel{}, AdminAccessLevel{}}
}

func (groupMember GroupMember) AccessLevelExists(accessLevel string) bool {
	accessLevels := groupMember.GetAccessLevels()
	exists := false
	for _, level := range accessLevels {
		if level.GetRole() == accessLevel {
			exists = true
		}
	}
	return exists
}

type AccessLevel interface {
	GetRole() string
}

type CasualAccessLevel struct{}

func (accessLevel CasualAccessLevel) GetRole() string {
	return "casual"
}

type ModeratorAccessLevel struct{}

func (accessLevel ModeratorAccessLevel) GetRole() string {
	return "moderator"
}

type AdminAccessLevel struct{}

func (accessLevel AdminAccessLevel) GetRole() string {
	return "admin"
}
