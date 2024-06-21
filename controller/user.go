package controller

import (
	models "GeekHub-backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterUserInput struct {
	Username string `gorm:"notnull" json:"username" binding:"required"`
	Email    string `gorm:"unique;notnull" json:"email" binding:"required"`
}

func RegisterUser(c *gin.Context) {
	var input RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	user.Username = input.Username
	user.Email = input.Email

	models.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{
		"created_user": user,
	})
}