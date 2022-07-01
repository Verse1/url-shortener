package main

import "github.com/gofiber/fiber/v2"

func initRoutes(app *fiber.App) {

	app.get("/:url", routes.connect)
	app.post("/api", routes.shorten)
}

func main() {
	app := fiber.New()
	initRoutes(app)
	app.Listen(":3000")
}