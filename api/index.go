package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/Verse1/url-shortener/api/routes"
	"github.com/joho/godotenv"
)

func initRoutes(app *fiber.App) {

	app.Get("/:url", routes.Connect)
	app.Post("/api", routes.Shorten)
}

func main() {

	env := godotenv.Load()
	if env != nil {
		fmt.Println(env)
	}


	app := fiber.New()
	initRoutes(app)
	app.Listen(":3000")
}