package database

import (
	"boilerplate/internal/domain"
	"time"

	"github.com/upper/db/v4"
)

const UsersTableName = "users"

type user struct {
	Id          uint64     `db:"id,omitempty"`
	Name        string     `db:"name"`
	Email       string     `db:"email"`
	Password    string     `db:"password"`
	Lat         float32    `db:"lat"`
	Lon         float32    `db:"lon"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type UserRepository interface {
	FindByEmail(email string) (domain.User, error)
	Save(user domain.User) (domain.User, error)
	FindById(id uint64) (domain.User, error)
	Update(user domain.User) (domain.User, error)
	Delete(id uint64) error
	GetCoordinates(user domain.User) (float32, float32, error)
	SetCoordinates(lat float32, lon float32, user domain.User) error
	GetUsersIdByArea(points map[string]map[string]float32) []uint64
}

type userRepository struct {
	coll db.Collection
}

func NewUserRepository(dbSession db.Session) UserRepository {
	return userRepository{
		coll: dbSession.Collection(UsersTableName),
	}
}

func (r userRepository) FindByEmail(email string) (domain.User, error) {
	var u user
	err := r.coll.Find(db.Cond{"email": email, "deleted_date": nil}).One(&u)
	if err != nil {
		return domain.User{}, err
	}

	return r.mapModelToDomain(u), nil
}

func (r userRepository) Save(user domain.User) (domain.User, error) {
	u := r.mapDomainToModel(user)
	u.CreatedDate, u.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&u)
	if err != nil {
		return domain.User{}, err
	}
	return r.mapModelToDomain(u), nil
}

func (r userRepository) FindById(id uint64) (domain.User, error) {
	var u user
	err := r.coll.Find(db.Cond{"id": id}).One(&u)
	if err != nil {
		return domain.User{}, err
	}
	return r.mapModelToDomain(u), nil
}

func (r userRepository) Update(user domain.User) (domain.User, error) {
	u := r.mapDomainToModel(user)
	u.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": u.Id}).Update(&u)
	if err != nil {
		return domain.User{}, err
	}
	return r.mapModelToDomain(u), nil
}

func (r userRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r userRepository) GetCoordinates(user domain.User) (float32, float32, error) {
	u := r.mapDomainToModel(user)
	u.UpdatedDate = time.Now()
	return u.Lat, u.Lon, nil
}

func (r userRepository) SetCoordinates(lat float32, lon float32, user domain.User) error {
	u := r.mapDomainToModel(user)
	u.UpdatedDate = time.Now()
	u.Lat = lat
	u.Lon = lon
	err := r.coll.Find(db.Cond{"id": u.Id}).Update(&u)
	if err != nil {
		return err
	}
	return nil
}

func (r userRepository) GetUsersIdByArea(points map[string]map[string]float32) []uint64 {
	var users []user
	query := r.coll.Find(db.Cond{"lat >": points["UpperLeftPoint"]["lat"], "lat <": points["BottomRightPoint"]["lat"], "lon <": points["UpperLeftPoint"]["lon"], "lon >": points["BottomRightPoint"]["lon"]})
	err := query.All(&users)
	if err != nil {
		return []uint64{}
	}
	usersId := make([]uint64, len(users))
	for i, user := range users {
		usersId[i] = user.Id
	}
	return usersId
}

func (r userRepository) mapDomainToModel(d domain.User) user {
	return user{
		Id:          d.Id,
		Name:        d.Name,
		Email:       d.Email,
		Password:    d.Password,
		Lat:         d.Lat,
		Lon:         d.Lon,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r userRepository) mapModelToDomain(m user) domain.User {
	return domain.User{
		Id:          m.Id,
		Name:        m.Name,
		Email:       m.Email,
		Password:    m.Password,
		Lat:         m.Lat,
		Lon:         m.Lon,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}
