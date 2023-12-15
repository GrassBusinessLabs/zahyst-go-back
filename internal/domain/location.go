package domain

import (
	"time"
)

type Location struct {
	Id          uint64
	UserId      uint64
	Type        string
	Address     string
	Lat         float64
	Lon         float64
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type Locations struct {
	Items []Location
	Total uint64
	Pages uint
}

func (loc Location) GetUserId() uint64 {
	return loc.UserId
}
