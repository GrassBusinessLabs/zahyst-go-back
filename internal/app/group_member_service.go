package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"log"
)

type GroupMemberService interface {
	AddGroupMember(accessCode string, userId uint64) (domain.GroupMember, error)
	ChangeAccessLevel(groupMember domain.GroupMember, newAccessLevel string) (domain.GroupMember, error)
	GetMembersList(p domain.Pagination, groupId uint64) (domain.GroupMembers, error)
	Find(id uint64) (interface{}, error)
	DeleteGroupMember(id uint64) error
}

type groupMemberService struct {
	groupMemberRepo database.GroupMemberRepository
	groupRepo       database.GroupRepository
}

func NewGroupMemberService(gmr database.GroupMemberRepository, gr database.GroupRepository) groupMemberService {
	return groupMemberService{
		groupMemberRepo: gmr,
		groupRepo:       gr,
	}
}

func (s groupMemberService) AddGroupMember(accessCode string, userId uint64) (domain.GroupMember, error) {
	grpMember, err := s.groupMemberRepo.AddGroupMember(accessCode, userId, s.groupRepo)
	if err != nil {
		log.Printf("GroupMemberService: %s", err)
		return domain.GroupMember{}, err
	}

	return grpMember, err
}

func (s groupMemberService) ChangeAccessLevel(groupMember domain.GroupMember, newAccessLevel string) (domain.GroupMember, error) {
	grpMember, err := s.groupMemberRepo.ChangeAccessLevel(groupMember, newAccessLevel)
	if err != nil {
		log.Printf("GroupMemberService: %s", err)
		return domain.GroupMember{}, err
	}

	return grpMember, err
}

func (s groupMemberService) GetMembersList(p domain.Pagination, groupId uint64) (domain.GroupMembers, error) {
	grpMembers, err := s.groupMemberRepo.GetMembersList(p, groupId)
	if err != nil {
		log.Printf("GroupMemberService: %s", err)
		return domain.GroupMembers{}, err
	}

	return grpMembers, err
}

func (s groupMemberService) DeleteGroupMember(id uint64) error {
	err := s.groupMemberRepo.DeleteGroupMember(id)
	if err != nil {
		log.Printf("GroupMemberService: %s", err)
		return err
	}

	return err
}

func (s groupMemberService) Find(id uint64) (interface{}, error) {
	group, err := s.groupMemberRepo.FindById(id)
	if err != nil {
		log.Printf("GroupMemberService: %s", err)
		return domain.GroupMember{}, err
	}

	return group, err
}
