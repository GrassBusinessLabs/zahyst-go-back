package requests

import (
	"boilerplate/internal/domain"
)

type CreateGroupRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateGroupRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (r CreateGroupRequest) ToDomainModel() (interface{}, error) {
	return domain.Group{
		Title:       r.Title,
		Description: r.Description,
	}, nil
}

func (r UpdateGroupRequest) ToDomainModel() (interface{}, error) {
	return domain.Group{
		Title:       r.Title,
		Description: r.Description,
	}, nil
}
