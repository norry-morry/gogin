package Handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"resume/entity"
	"resume/repository"
)

func GetUsers(c *gin.Context) {
	var user []entity.User
	err := repository.GetAllUsers(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func GetUserById(c *gin.Context) {
	id := c.Params.ByName("id")
	var user entity.User
	err := repository.GetUserById(&user, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func CreateUser(c *gin.Context) {
	var user entity.User
	jsonErr := c.BindJSON(&user)
	if jsonErr != nil {
		return
	}
	err := repository.CreateUser(&user)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusCreated, user)
	}
}

func UpdateUser(c *gin.Context) {
	var user entity.User
	id := c.Params.ByName("id")
	err := repository.GetUserById(&user, id)
	if err != nil {
		c.JSON(http.StatusNotFound, user)
	}
	c.BindJSON(&user)
	err = repository.UpdateUser(&user, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(c *gin.Context) {
	var user entity.User
	id := c.Params.ByName("id")
	err := repository.DeleteUser(&user, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, gin.H{"id" + id: "is deleted"})
	}
}
