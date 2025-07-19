package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/bars-be/controllers"
	"github.com/heronhoga/bars-be/middlewares"
)

func AuthRoutes(app *fiber.App) {
	app.Post("/register", middlewares.CheckAppKey, controllers.Register)
	app.Post("/login", middlewares.CheckAppKey, controllers.Login)
}