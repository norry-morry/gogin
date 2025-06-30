package user

type Repository interface {
	FindAll() ([]User, error)
	FindById(id uint) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}
