package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Balance   float32
}

type UserDTO struct {
	ID        uint    `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Balance   float32 `json:"balance"`
}

type UserResponse struct {
	ID        uint    `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Balance   float32 `json:"balance"`
}
