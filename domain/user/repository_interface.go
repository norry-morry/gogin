package user

// Repository はユーザー情報へのリポジトリ操作を定義するインターフェースです。
type Repository interface {
	FindAll() ([]User, error)
	FindByID(id uint) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}
