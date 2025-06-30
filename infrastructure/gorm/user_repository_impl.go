package gorm

import (
	"gorm.io/gorm"
	"resume/domain/user"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository は Repository インターフェースの GORM 実装を返します。
func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db}
}

func (r *userRepository) FindAll() ([]user.User, error) {
	var users []user.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) FindById(id uint) (*user.User, error) {
	var u user.User
	err := r.db.First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

func (r *userRepository) Update(u *user.User) error {
	return r.db.Save(u).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&user.User{}, id).Error
}
