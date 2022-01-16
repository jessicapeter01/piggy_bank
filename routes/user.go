package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jessicapeter01/piggy_bank/app/controllers"
)

func UserRoute(route fiber.Router) {
	route.Get("", controllers.GetUsers)
	route.Get("/:id", controllers.GetUser)
	route.Post("", controllers.CreateUser)
	route.Delete("/:id", controllers.DeleteUser)
}
