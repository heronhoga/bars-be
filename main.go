package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/bars-be/config"
	"github.com/heronhoga/bars-be/routes"
	"github.com/heronhoga/bars-be/utils"
)

func main() {
	app := fiber.New()

	//load env
	utils.LoadEnv()

	//connect to database
	config.InitDB()

	//routes config
	routes.AuthRoutes(app)

	app.Listen(":3000")
}