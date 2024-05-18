package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	log.Println("🚀 Starting server...")

	// Create a new Fiber instance
	app := fiber.New()
	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// Start the server
	log.Println("🚀 Server is running on port", config.PORT)
	log.Fatal(app.Listen(fmt.Sprintf(":%v", config.PORT)))
}
