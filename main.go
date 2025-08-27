package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/heronhoga/bars-be/config"
	"github.com/heronhoga/bars-be/routes"
	"github.com/heronhoga/bars-be/utils"
)

func main() {
	//load env
	utils.LoadEnv()
	frontEndApp := os.Getenv("FRONTEND_APP")
	fmt.Println(frontEndApp)
	app := fiber.New(
		fiber.Config{
			BodyLimit: 6 * 1024 * 1024,
		},
	)

	//cors config
	app.Use(cors.New(cors.Config{
        AllowOrigins: frontEndApp,
        AllowHeaders: "Origin, Content-Type, Accept, app-key",
        AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
        AllowCredentials: true,
        }))
		


	//connect to database
	config.InitDB()

	//routes config
	routes.AuthRoutes(app)
	routes.BeatRoutes(app)
	routes.LikesRoutes(app)
	routes.ProfileRoutes(app)

	app.Listen(":8000")
}