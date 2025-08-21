package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/bars-be/controllers"
	"github.com/heronhoga/bars-be/middlewares"
)

func ProfileRoutes(app *fiber.App) {
	app.Get("/profile", middlewares.CheckAppKey, middlewares.CheckJWT, controllers.GetProfile)
	app.Get("/beatbyuser", middlewares.CheckAppKey, middlewares.CheckJWT, controllers.GetBeatByUser)
	app.Get("/likedbyuser", middlewares.CheckAppKey, middlewares.CheckJWT, controllers.GetLikedBeatByUser)
}