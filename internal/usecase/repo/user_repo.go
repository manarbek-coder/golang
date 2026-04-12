package repo

import (
	"practice-7/internal/entity"
	"practice-7/pkg/postgres"

	"github.com/google/uuid"
)

type UserRepo struct {
	pg *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg: pg}
}

func (r *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	r.pg.DB.Create(user)
	return user, nil
}

func (r *UserRepo) LoginUser(username string) (*entity.User, error) {
	var user entity.User
	r.pg.DB.Where("username = ?", username).First(&user)
	return &user, nil
}

func (r *UserRepo) GetByID(id string) (*entity.User, error) {
	var user entity.User
	uid, _ := uuid.Parse(id)
	r.pg.DB.First(&user, "id = ?", uid)
	return &user, nil
}

func (r *UserRepo) PromoteUser(id string) error {
	return r.pg.DB.Model(&entity.User{}).
		Where("id = ?", id).
		Update("role", "admin").Error
}
