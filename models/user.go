package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name string `json:"name"`
	Age  uint8  `json:"age"`
}
