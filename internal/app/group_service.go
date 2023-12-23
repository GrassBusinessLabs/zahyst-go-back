package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"log"
	"math/rand"
	"time"
)

type GroupService interface {
	Save(group domain.Group) (domain.Group, error)
	Update(group domain.Group) (domain.Group, error)
	Delete(id uint64) error
	Find(uint64) (interface{}, error)
	GetList(p domain.Pagination) (domain.Groups, error)
	GetAccessCode(group domain.Group) string
}

type groupService struct {
	groupRepo database.GroupRepository
}

func NewGroupService(gr database.GroupRepository) groupService {
	return groupService{
		groupRepo: gr,
	}
}

func (s groupService) Save(group domain.Group) (domain.Group, error) {
	group.AcessCode = s.GenerateAccessCode()
	grp, err := s.groupRepo.Save(group)
	if err != nil {
		log.Printf("GroupService: %s", err)
		return domain.Group{}, err
	}

	return grp, err
}

func (s groupService) Update(group domain.Group) (domain.Group, error) {
	grp, err := s.groupRepo.Update(group)
	if err != nil {
		log.Printf("GroupService: %s", err)
		return domain.Group{}, err
	}

	return grp, err
}

func (s groupService) Delete(id uint64) error {
	err := s.groupRepo.Delete(id)
	if err != nil {
		log.Printf("GroupService: %s", err)
		return err
	}

	return nil
}

func (s groupService) GetList(p domain.Pagination) (domain.Groups, error) {
	groups, err := s.groupRepo.GetList(p)
	if err != nil {
		log.Printf("GroupService: %s", err)
		return domain.Groups{}, err
	}

	return groups, err
}

func (s groupService) Find(id uint64) (interface{}, error) {
	group, err := s.groupRepo.FindById(id)
	if err != nil {
		log.Printf("GroupService: %s", err)
		return domain.Group{}, err
	}

	return group, err
}

func (s groupService) GetAccessCode(group domain.Group) string {
	accessCode := s.groupRepo.GetAccessCode(group)
	return accessCode
}

func (s groupService) GenerateAccessCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLength := 6
	id := make([]byte, idLength)
	for i := range id {
		id[i] = chars[rnd.Intn(len(chars))]
	}
	return string(id)
}
