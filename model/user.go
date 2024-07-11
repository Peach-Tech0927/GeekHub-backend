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