package repos

import (
	"errors"

	"github.com/jessicapeter01/piggy_bank/app/models"
	"github.com/jessicapeter01/piggy_bank/database"
	"gorm.io/gorm"
)

func GetAllUsers() []models.User {
	db := database.DBConn
	var users []models.User
	db.Find(&users)
	return users
}

func GetUserByID(UserId uint) (models.User, error) {
	db := database.DBConn
	var User models.User
	err := db.First(&User, UserId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User, errors.New("user Not Found")
	}
	return User, nil
}

func CreateUser(firstName string, lastName string, balance float32) (models.User, error) {
	db := database.DBConn

	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Balance:   balance,
	}

	if err := db.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func DeleteUser(itemId uint) error {
	db := database.DBConn

	var item models.User
	err := db.First(&item, itemId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	}
	return db.Delete(&item, itemId).Error
}
