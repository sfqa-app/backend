package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//	@Description	User account
type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null;default:null"`
	Username string `json:"username" gorm:"unique; not null;default:null"`
	Email    string `json:"email" gorm:"unique; not null;default:null"`
	Password string `json:"password" gorm:"not null;default:null"`
}

func (user *User) EncryptPassword(password string) error {
	bytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return nil
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
