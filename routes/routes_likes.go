package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/bars-be/controllers"
	"github.com/heronhoga/bars-be/middlewares"
)

func LikesRoutes(app *fiber.App) {
	app.Post("/like", middlewares.CheckAppKey, middlewares.CheckJWT, controllers.Like)
}