package repository

import (
	"fmt"
	"resume/entity"
	"resume/library/database"
)

func GetAllUsers(user *[]entity.User) (err error) {
	if err = database.DB.Find(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserById(user *entity.User, id string) (err error) {
	if err = database.DB.Where("id=?", id).First(user).Error; err != nil {
		return err
	}
	return nil
}

func CreateUser(user *entity.User) (err error) {
	if err = database.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func UpdateUser(user *entity.User, id string) (err error) {
	fmt.Println(id)
	fmt.Println(user)
	if err = database.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUser(user *entity.User, id string) (err error) {
	if err = database.DB.Where("id=?", id).Delete(user).Error; err != nil {
		return err
	}
	return nil
}
