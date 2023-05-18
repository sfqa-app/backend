package models

import "gorm.io/gorm"

//	@Description	User account info
type UserInfo struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

//	@Description	User account
type User struct {
	gorm.Model
	UserInfo
}
