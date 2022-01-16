package models

import (
	"time"

	"gorm.io/gorm"
)

type Goal struct {
	gorm.Model
	Title   string    `json:"title"`
	EndDate time.Time `json:"endDate"`
	Total   float32   `json:"total"`
	Saved   float32   `json:"saved"`
	User    User      `json:"user"`
	UserID  uint      `json:"userId"`
}

type GoalDTO struct {
	ID      uint      `json:"id"`
	Title   string    `json:"title"`
	EndDate time.Time `json:"endDate"`
	Total   float32   `json:"total"`
	Saved   float32   `json:"saved"`
	UserID  uint      `json:"userId"`
}

type GoalResponse struct {
	ID      uint      `json:"id"`
	Title   string    `json:"title"`
	EndDate time.Time `json:"endDate"`
	Total   float32   `json:"total"`
	Saved   float32   `json:"saved"`
	UserID  uint      `json:"userId"`
}
