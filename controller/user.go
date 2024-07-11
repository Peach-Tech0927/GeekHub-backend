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

	var user models.User
	user.Username = input.Username
	user.Email = input.Email

	if err := user.Create(); err != nil {
		if err == models.ErrEmailInUse {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
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

func GetUserInfo(c *gin.Context) {
	userId, err := token.ExtractTokenId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	user, err := models.GetUserInfoById(userId)
	if err != nil {
		if err == models.ErrFailedToGetUserInfo {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrUnexpected.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func SignIn(c *gin.Context) {
	var input SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	user, err := models.IdentifyUserByEmail(input.Email)
	if err != nil {
		if err == models.ErrFailedToGetUserInfo {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrUnexpected.Error()})
		}
		return
	}

	token, err := token.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}