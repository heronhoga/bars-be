package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/bars-be/controllers"
	"github.com/heronhoga/bars-be/middlewares"
)

func BeatRoutes(app *fiber.App) {
	//authorized users
	app.Post("/beat", middlewares.CheckAppKey, middlewares.CheckJWT, controllers.CreateNewBeat)
	app.Get("/beat", middlewares.CheckAppKey, middlewares.CheckJWT, controllers.GetAllBeats)
	app.Delete("/beat", middlewares.CheckAppKey, middlewares.CheckJWT, controllers.DeleteBeat)
	app.Put("/beat/:beatid", middlewares.CheckAppKey, middlewares.CheckJWT, controllers.EditBeat)

	//unauthorized users
	app.Get("/favoritebeats", middlewares.CheckAppKey, controllers.GetFavoriteBeats)
}