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

func GetGoal(c *fiber.Ctx) error {
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

	goal, err := repos.GetGoalByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Goal not found",
		})
	}

	result := &models.GoalResponse{}
	if err := copier.Copy(&result, &goal); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"goal": result,
		},
	})
}

func GetGoals(c *fiber.Ctx) error {
	goals := repos.GetAllGoals()
	result := []models.GoalResponse{}
	if err := copier.Copy(&result, &goals); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"goal": result,
		},
	})
}

// Create a risk
func CreateGoal(c *fiber.Ctx) error {
	body := models.GoalDTO{}

	if err := utils.ParseBodyAndValidate(c, &body); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}

	goal, err := repos.CreateGoal(body.Title, body.EndDate, body.Total, body.Saved, body.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot create goal",
		})
	}

	result := &models.GoalResponse{}
	if err = copier.Copy(&result, &goal); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"goal": result,
		},
	})
}

func DeleteGoal(c *fiber.Ctx) error {
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

	// find Goal and return
	if err := repos.DeleteGoal(id); err == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}

	// if Goal not available
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Goal not found",
	})
}
