package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/bars-be/controllers"
	"github.com/heronhoga/bars-be/middlewares"
)

func BeatRoutes(app *fiber.App) {
	app.Post("/beat", middlewares.CheckAppKey, middlewares.CheckJWT, controllers.CreateNewBeat)
}