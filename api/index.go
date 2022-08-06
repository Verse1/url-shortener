package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware"
	// "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func initRoutes(app *fiber.App) {

	app.get("/:url", routes.connect)
	app.post("/api", routes.shorten)
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