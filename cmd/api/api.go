package api

import (
	"github.com/gofiber/fiber/v2"
)

func Initialize() {
	app := fiber.New()

	RegisterRoutes(app)

	app.Listen(":3000")
}
