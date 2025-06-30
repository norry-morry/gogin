package user

import "time"

type User struct {
	ID          uint64    `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email" gorm:"uniqueIndex"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	//DeletedAt 	time.Time	`json:"deleted_at"`
}
