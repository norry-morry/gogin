package user

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"resume/domain/user"
	"testing"
)

type mockUserRepository struct {
	FindAllFn  func() ([]user.User, error)
	FindByIdFn func(uint) (*user.User, error)
	CreateFn   func(*user.User) error
	UpdateFn   func(*user.User) error
	DeleteFn   func(uint) error
}

func (m *mockUserRepository) FindAll() ([]user.User, error) {
	return m.FindAllFn()
}
func (m *mockUserRepository) FindById(id uint) (*user.User, error) {
	return m.FindByIdFn(id)
}
func (m *mockUserRepository) Create(user *user.User) error {
	return m.CreateFn(user)
}
func (m *mockUserRepository) Update(user *user.User) error {
	return m.UpdateFn(user)
}
func (m *mockUserRepository) Delete(id uint) error {
	return m.DeleteFn(id)
}

func TestFindAllUsers_Success(t *testing.T) {
	mock := &mockUserRepository{
		FindAllFn: func() ([]user.User, error) {
			return []user.User{
				{ID: 1, Name: "test1"},
				{ID: 2, Name: "test2"},
			}, nil
		},
	}

	usecase := NewUserUsecase(mock)
	users, err := usecase.GetAllUsers()
	assert.NoErrorf(t, err, "Get all users")
	assert.Len(t, users, 2)
	assert.Equal(t, users[0].Name, "test1")
}

func TestFindAllUsers_Failure(t *testing.T) {
	mock := &mockUserRepository{
		FindAllFn: func() ([]user.User, error) {
			return nil, errors.New("DB error")
		},
	}

	usecase := NewUserUsecase(mock)

	users, err := usecase.GetAllUsers()

	assert.Error(t, err)
	assert.Nil(t, users)
}
