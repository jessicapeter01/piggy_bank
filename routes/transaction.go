package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jessicapeter01/piggy_bank/app/controllers"
)

func TransactionRoute(route fiber.Router) {
	route.Get("", controllers.GetTransactions)
	route.Get("/:id", controllers.GetTransaction)
	route.Post("", controllers.CreateTransaction)
	route.Delete("/:id", controllers.DeleteTransaction)
}
