package database

import (
	"boilerplate/internal/domain"
	"math"
	"time"

	"github.com/upper/db/v4"
)

const GroupsTableName = "groups"

type group struct {
	Id          uint64     `db:"id,omitempty"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	UserId      uint64     `db:"user_id"`
	AccessCode  string     `db:"access_code"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type GroupRepository interface {
	Save(group domain.Group) (domain.Group, error)
	Update(group domain.Group) (domain.Group, error)
	Delete(id uint64) error
	FindById(id uint64) (domain.Group, error)
	GetList(p domain.Pagination) (domain.Groups, error)
	GetAccessCode(group domain.Group) string
	GetGroupByAccessCode(accessCode string) (domain.Group, error)
}

type groupRepository struct {
	coll db.Collection
}

func NewGroupRepository(dbSession db.Session) groupRepository {
	return groupRepository{
		coll: dbSession.Collection(GroupsTableName),
	}
}

func (r groupRepository) Save(group domain.Group) (domain.Group, error) {
	grp := r.mapDomainToModel(group)
	grp.CreatedDate, grp.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&grp)
	if err != nil {
		return domain.Group{}, err
	}
	return r.mapModelToDomain(grp), nil
}

func (r groupRepository) Update(group domain.Group) (domain.Group, error) {
	grp := r.mapDomainToModel(group)
	grp.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": grp.Id}).Update(&grp)
	if err != nil {
		return domain.Group{}, err
	}
	return r.mapModelToDomain(grp), nil
}

func (r groupRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r groupRepository) FindById(id uint64) (domain.Group, error) {
	var grp group
	err := r.coll.Find(db.Cond{"id": id}).One(&grp)
	if err != nil {
		return domain.Group{}, err
	}
	return r.mapModelToDomain(grp), nil
}

func (r groupRepository) GetList(p domain.Pagination) (domain.Groups, error) {
	var data []group
	query := r.coll.Find(db.Cond{})
	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.Groups{}, err
	}

	groups := r.mapModelToDomainPagination(data)

	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Groups{}, err
	}

	groups.Total = totalCount
	groups.Pages = uint(math.Ceil(float64(groups.Total) / float64(p.CountPerPage)))

	return groups, nil
}

func (r groupRepository) GetAccessCode(group domain.Group) string {
	return group.AcessCode
}

func (r groupRepository) GetGroupByAccessCode(accessCode string) (domain.Group, error) {
	var grp group
	err := r.coll.Find(db.Cond{"access_code": accessCode}).One(&grp)
	if err != nil {
		return domain.Group{}, err
	}
	return r.mapModelToDomain(grp), nil
}

func (r groupRepository) mapDomainToModel(d domain.Group) group {
	return group{
		Id:          d.Id,
		UserId:      d.UserId,
		Title:       d.Title,
		Description: d.Description,
		AccessCode:  d.AcessCode,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r groupRepository) mapModelToDomain(m group) domain.Group {
	return domain.Group{
		Id:          m.Id,
		UserId:      m.UserId,
		Title:       m.Title,
		Description: m.Description,
		AcessCode:   m.AccessCode,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (f groupRepository) mapModelToDomainPagination(groups []group) domain.Groups {
	new_groups := make([]domain.Group, len(groups))
	for i, group := range groups {
		new_groups[i] = f.mapModelToDomain(group)
	}
	return domain.Groups{Items: new_groups}
}
