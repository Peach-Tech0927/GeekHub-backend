package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
}

var (
	ErrEmailInUse = errors.New("email is already in use")
	ErrFailedToGetUserInfo = errors.New("failed to get user info")
	ErrUnexpected = errors.New("an unexpected error occurred")
)

func (u *User) Create() error {
	var existingUser User
	if err := DB.Where("email = ?", u.Email).First(&existingUser).Error; err == nil {
		return ErrEmailInUse
	}

	if err := DB.Create(u).Error; err != nil {
		return err
	}

	return nil
}

func GetUserInfoById(id uint) (*User, error) {
	var user User
	if err := DB.Where("ID = ?", id).First(&user).Error; err != nil {
		return nil, ErrFailedToGetUserInfo
	}

	return &user, nil
}