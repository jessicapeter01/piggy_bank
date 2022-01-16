package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jessicapeter01/piggy_bank/app/models"
	"github.com/jessicapeter01/piggy_bank/app/repos"
	"github.com/jessicapeter01/piggy_bank/app/utils"
	"github.com/jinzhu/copier"
)

func GetTransaction(c *fiber.Ctx) error {
	// get parameter value
	paramId := c.Params("id")
	var id uint
	// convert parameter value string to int
	if v, err := strconv.ParseUint(paramId, 10, 32); err == nil {
		id = uint(v)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse ID",
		})
	}

	transaction, err := repos.GetTransactionByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Transaction not found",
		})
	}

	result := &models.TransactionResponse{}
	if err := copier.Copy(&result, &transaction); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"transaction": result,
		},
	})
}

func GetTransactions(c *fiber.Ctx) error {
	transactions := repos.GetAllTransactions()
	result := []models.TransactionResponse{}
	if err := copier.Copy(&result, &transactions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"transaction": result,
		},
	})
}

// Create a risk
func CreateTransaction(c *fiber.Ctx) error {
	body := models.TransactionDTO{}

	if err := utils.ParseBodyAndValidate(c, &body); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}

	transaction, err := repos.CreateTransaction(body.Date, body.GoalID, body.Total)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot create transaction",
		})
	}

	result := &models.TransactionResponse{}
	if err = copier.Copy(&result, &transaction); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"transaction": result,
		},
	})
}

func DeleteTransaction(c *fiber.Ctx) error {
	// get parameter value
	paramId := c.Params("id")
	var id uint
	// convert parameter value string to int
	if v, err := strconv.ParseUint(paramId, 10, 32); err == nil {
		id = uint(v)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse ID",
		})
	}

	// find Transaction and return
	if err := repos.DeleteTransaction(id); err == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}

	// if Transaction not available
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Transaction not found",
	})
}
