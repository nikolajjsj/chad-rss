package route

import (
	"chadrss/api/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Authentication
	api.Post("/signup", controller.Signup)
	api.Get("/login", controller.Login)

}
