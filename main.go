package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/heronhoga/bars-be/config"
	"github.com/heronhoga/bars-be/routes"
	"github.com/heronhoga/bars-be/utils"
)

func main() {
	app := fiber.New(
		fiber.Config{
			BodyLimit: 6 * 1024 * 1024,
		},
	)

	//cors config
	app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:3000",
        AllowHeaders: "Origin, Content-Type, Accept, app-key",
        AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
        AllowCredentials: true,
        }))
		
	//load env
	utils.LoadEnv()

	//connect to database
	config.InitDB()

	//routes config
	routes.AuthRoutes(app)
	routes.BeatRoutes(app)
	routes.LikesRoutes(app)
	routes.ProfileRoutes(app)

	app.Listen(":8000")
}