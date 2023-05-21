package models

import "gorm.io/gorm"

//	@Description	User account info
type UserInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//	@Description	User account
type User struct {
	gorm.Model
	UserInfo
}
