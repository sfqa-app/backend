package models

import (
	valid "github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// @Description	User account
type User struct {
	gorm.Model
	Name            string `json:"name"`
  Username        string `json:"username" gorm:"unique;default:null"`
	Email           string `json:"email" gorm:"unique; not null;default:null"`
	Password        string `json:"-" gorm:"not null;default:null"`
	IsEmailVerified bool   `json:"is_email_verified" gorm:"default:false"`
}

func NewUser(email, password string) *User {
	return &User{
		Email:    email,
		Password: password,
	}
}

func (user *User) EncryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return nil
}

func (user *User) IsPasswordMatch(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) IsValidEmail() bool {
	return valid.IsEmail(user.Email)
}
