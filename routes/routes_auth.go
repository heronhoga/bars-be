package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/bars-be/controllers"
)

func AuthRoutes(app *fiber.App) {
	app.Post("/register", controllers.Register)
}