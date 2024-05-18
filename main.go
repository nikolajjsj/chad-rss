package main

import (
	"chadrss/api/route"
	"chadrss/config"
	"chadrss/database"
	"chadrss/frontend"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	log.Println("🚀 Starting server...")

	// Connect to the database
	database.ConnectDB()

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
	// Serve static files
	app.All("/*", filesystem.New(filesystem.Config{
		Root:         frontend.Dist(),
		NotFoundFile: "index.html",
		Index:        "index.html",
	}))
	route.SetupRoutes(app)

	// Start the server
	log.Println("🚀 Server is running on port", config.PORT)
	log.Fatal(app.Listen(fmt.Sprintf(":%v", config.PORT)))
}
