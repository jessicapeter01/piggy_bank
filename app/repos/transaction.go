package repos

import (
	"errors"
	"time"

	"github.com/jessicapeter01/piggy_bank/app/models"
	"github.com/jessicapeter01/piggy_bank/database"
	"gorm.io/gorm"
)

func GetAllTransactions() []models.Transaction {
	db := database.DBConn
	var transactions []models.Transaction
	db.Find(&transactions)
	return transactions
}

func GetTransactionByID(transactionId uint) (models.Transaction, error) {
	db := database.DBConn
	var transaction models.Transaction
	err := db.First(&transaction, transactionId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return transaction, errors.New("transaction Not Found")
	}
	return transaction, nil
}

func CreateTransaction(date time.Time, goalID uint, total float32) (models.Transaction, error) {
	db := database.DBConn

	transaction := models.Transaction{
		Date:   date,
		GoalID: goalID,
		Total:  total,
	}

	if err := db.Create(&transaction).Error; err != nil {
		return models.Transaction{}, err
	}
	return transaction, nil
}

func DeleteTransaction(itemId uint) error {
	db := database.DBConn

	var item models.Transaction
	err := db.First(&item, itemId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("transaction not found")
	}
	return db.Delete(&item, itemId).Error
}
