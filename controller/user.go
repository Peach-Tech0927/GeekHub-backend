package controller

import (
	models "GeekHub-backend/model"
	"GeekHub-backend/utils/token"
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

	// emailが使われている(err == nil)ときにエラーを返す
	var existingUser models.User
	if err := models.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already in use"})
		return
	}

	var user models.User
	user.Username = input.Username
	user.Email = input.Email

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := token.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"created_user": user, "token": token,
	})
}