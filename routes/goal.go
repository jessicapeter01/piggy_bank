package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jessicapeter01/piggy_bank/app/controllers"
)

func GoalRoute(route fiber.Router) {
	route.Get("", controllers.GetGoals)
	route.Get("/:id", controllers.GetGoal)
	route.Post("", controllers.CreateGoal)
	route.Delete("/:id", controllers.DeleteGoal)
}
