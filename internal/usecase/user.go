package usecase

import (
	"Assignment5/internal/repository/users"
	"Assignment5/pkg/modules"
)

type UserUsecase struct {
	repo *users.Repository
}

func NewUserUsecase(r *users.Repository) *UserUsecase {
	return &UserUsecase{r}
}

func (u *UserUsecase) GetUsers(f modules.UserFilter) (modules.PaginatedResponse, error) {
	return u.repo.GetPaginatedUsers(f)
}

func (u *UserUsecase) GetCommonFriends(a, b int) ([]modules.User, error) {
	return u.repo.GetCommonFriends(a, b)
}
