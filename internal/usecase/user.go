package usecase

import (
	"fmt"
	"practice-7/internal/entity"
	"practice-7/internal/usecase/repo"
	"practice-7/utils"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	user, err := u.repo.RegisterUser(user)
	if err != nil {
		return nil, "", err
	}
	return user, "session", nil
}

func (u *UserUseCase) LoginUser(input *entity.LoginUserDTO) (string, error) {
	user, err := u.repo.LoginUser(input.Username)
	if err != nil {
		return "", err
	}
	if !utils.CheckPassword(user.Password, input.Password) {
		return "", fmt.Errorf("wrong password")
	}
	return utils.GenerateJWT(user.ID, user.Role)
}

func (u *UserUseCase) GetMe(userID string) (*entity.User, error) {
	return u.repo.GetByID(userID)
}

func (u *UserUseCase) PromoteUser(id string) error {
	return u.repo.PromoteUser(id)
}
