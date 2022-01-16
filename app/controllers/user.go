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

func GetUser(c *fiber.Ctx) error {
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

	user, err := repos.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	result := &models.UserResponse{}
	if err := copier.Copy(&result, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": result,
		},
	})
}

func GetUsers(c *fiber.Ctx) error {
	users := repos.GetAllUsers()
	result := []models.UserResponse{}
	if err := copier.Copy(&result, &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": result,
		},
	})
}

// Create a risk
func CreateUser(c *fiber.Ctx) error {
	body := models.UserDTO{}

	if err := utils.ParseBodyAndValidate(c, &body); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}

	user, err := repos.CreateUser(body.FirstName, body.LastName, body.Balance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot create user",
		})
	}

	result := &models.UserResponse{}
	if err = copier.Copy(&result, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": result,
		},
	})
}

func DeleteUser(c *fiber.Ctx) error {
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

	// find User and return
	if err := repos.DeleteUser(id); err == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}

	// if User not available
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "User not found",
	})
}
