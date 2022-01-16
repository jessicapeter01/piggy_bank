package repos

import (
	"errors"
	"time"

	"github.com/jessicapeter01/piggy_bank/app/models"
	"github.com/jessicapeter01/piggy_bank/database"
	"gorm.io/gorm"
)

func GetAllGoals() []models.Goal {
	db := database.DBConn
	var goals []models.Goal
	db.Find(&goals)
	return goals
}

func GetGoalByID(goalID uint) (models.Goal, error) {
	db := database.DBConn
	var goal models.Goal
	err := db.First(&goal, goalID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return goal, errors.New("goal Not Found")
	}
	return goal, nil
}

func CreateGoal(title string, endDate time.Time, total float32, saved float32, userID uint) (models.Goal, error) {
	db := database.DBConn

	goal := models.Goal{
		Title:   title,
		EndDate: endDate,
		Total:   total,
		Saved:   saved,
		UserID:  userID,
	}

	if err := db.Create(&goal).Error; err != nil {
		return models.Goal{}, err
	}
	return goal, nil
}

func DeleteGoal(itemId uint) error {
	db := database.DBConn

	var item models.Goal
	err := db.First(&item, itemId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("goal not found")
	}
	return db.Delete(&item, itemId).Error
}
