package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Date   time.Time `json:"date"`
	Goal   Goal      `json:"goal"`
	GoalID uint      `json:"goalId"`
	Total  float32   `json:"total"`
}

type TransactionDTO struct {
	ID     uint      `json:"id"`
	Date   time.Time `json:"date"`
	GoalID uint      `json:"goalId"`
	Total  float32   `json:"total"`
}

type TransactionResponse struct {
	ID     uint      `json:"id"`
	Date   time.Time `json:"date"`
	GoalID uint      `json:"goalId"`
	Total  float32   `json:"total"`
}
