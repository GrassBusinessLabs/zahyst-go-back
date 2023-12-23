package database

import (
	"boilerplate/internal/domain"
	"math"
	"time"

	"github.com/upper/db/v4"
)

const GroupMembersTableName = "group_members"

type groupMember struct {
	Id          uint64     `db:"id,omitempty"`
	UserId      uint64     `db:"user_id"`
	GroupId     uint64     `db:"group_id"`
	AccessLevel string     `db:"access_level"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type GroupMemberRepository interface {
	AddGroupMember(accessCode string, userId uint64, gr GroupRepository) (domain.GroupMember, error)
	ChangeAccessLevel(groupMember domain.GroupMember, newAccessLevel string) (domain.GroupMember, error)
	GetMembersList(p domain.Pagination, groupId uint64) (domain.GroupMembers, error)
	FindById(id uint64) (domain.GroupMember, error)
	DeleteGroupMember(id uint64) error
}

type groupMemberRepository struct {
	coll db.Collection
}

func NewGroupMemberRepository(dbSession db.Session) groupMemberRepository {
	return groupMemberRepository{
		coll: dbSession.Collection(GroupMembersTableName),
	}
}

func (r groupMemberRepository) AddGroupMember(accessCode string, userId uint64, gr GroupRepository) (domain.GroupMember, error) {
	grp, err := gr.GetGroupByAccessCode(accessCode)
	if err != nil {
		return domain.GroupMember{}, err
	}
	var grpMember groupMember
	grpMember.GroupId = grp.Id
	grpMember.UserId = userId
	grpMember.AccessLevel = "casual"
	grpMember.CreatedDate, grpMember.UpdatedDate = time.Now(), time.Now()
	err = r.coll.InsertReturning(&grpMember)
	if err != nil {
		return domain.GroupMember{}, err
	}
	return r.mapModelToDomain(grpMember), nil
}

func (r groupMemberRepository) ChangeAccessLevel(groupMember domain.GroupMember, newAccessLevel string) (domain.GroupMember, error) {
	grpMember := r.mapDomainToModel(groupMember)
	grpMember.AccessLevel = newAccessLevel
	err := r.coll.Find(db.Cond{"id": grpMember.Id}).Update(&grpMember)
	if err != nil {
		return domain.GroupMember{}, err
	}
	return r.mapModelToDomain(grpMember), nil
}

func (r groupMemberRepository) GetMembersList(p domain.Pagination, groupId uint64) (domain.GroupMembers, error) {
	var data []groupMember
	query := r.coll.Find(db.Cond{"group_id": groupId})
	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.GroupMembers{}, err
	}

	groupMembers := r.mapModelToDomainPagination(data)

	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.GroupMembers{}, err
	}

	groupMembers.Total = totalCount
	groupMembers.Pages = uint(math.Ceil(float64(groupMembers.Total) / float64(p.CountPerPage)))

	return groupMembers, nil
}

func (r groupMemberRepository) DeleteGroupMember(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r groupMemberRepository) FindById(id uint64) (domain.GroupMember, error) {
	var grpMember groupMember
	err := r.coll.Find(db.Cond{"id": id}).One(&grpMember)
	if err != nil {
		return domain.GroupMember{}, err
	}
	return r.mapModelToDomain(grpMember), nil
}

func (r groupMemberRepository) mapDomainToModel(d domain.GroupMember) groupMember {
	return groupMember{
		Id:          d.Id,
		UserId:      d.UserId,
		GroupId:     d.GroupId,
		AccessLevel: d.AccessLevel,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r groupMemberRepository) mapModelToDomain(m groupMember) domain.GroupMember {
	return domain.GroupMember{
		Id:          m.Id,
		UserId:      m.UserId,
		GroupId:     m.GroupId,
		AccessLevel: m.AccessLevel,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (f groupMemberRepository) mapModelToDomainPagination(groupMembers []groupMember) domain.GroupMembers {
	new_group_members := make([]domain.GroupMember, len(groupMembers))
	for i, group_member := range groupMembers {
		new_group_members[i] = f.mapModelToDomain(group_member)
	}
	return domain.GroupMembers{Items: new_group_members}
}
