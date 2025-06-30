package user

import "resume/domain/user"

type UserUsecase interface {
	GetAllUsers() ([]user.User, error)
	GetUserById(id uint) (*user.User, error)
	CreateUser(name, displayName string, email string) error
	UpdateUser(id uint, name, displayName string, email string) error
	DeleteUser(id uint) error
}

type userUsecase struct {
	repo user.Repository
}

func (uc *userUsecase) GetAllUsers() ([]user.User, error) {
	return uc.repo.FindAll()
}

func (uc userUsecase) GetUserById(id uint) (*user.User, error) {
	return uc.repo.FindById(id)
}

func (uc userUsecase) CreateUser(name, displayName string, email string) error {
	u := &user.User{
		Name:        name,
		DisplayName: displayName,
		Email:       email,
	}
	return uc.repo.Create(u)
}

func (uc userUsecase) UpdateUser(id uint, name, displayName string, email string) error {
	u, err := uc.repo.FindById(id)
	if err != nil {
		return err
	}
	u.Name = name
	u.DisplayName = displayName
	u.Email = email
	return uc.repo.Update(u)
}

func (uc userUsecase) DeleteUser(id uint) error {
	_, err := uc.repo.FindById(id)
	if err != nil {
		return err
	}
	return uc.repo.Delete(id)
}

func NewUserUsecase(repo user.Repository) UserUsecase {
	return &userUsecase{repo: repo}
}
