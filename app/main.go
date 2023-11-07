package main

import (
	"face-recognition-golang/db"
	"face-recognition-golang/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	r := app.Group("/face")
	// Initialize default config (Assign the middleware to /metrics)
	r.Get("/metrics", monitor.New())
	// Health check
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// Initialize MongoDB
	mongoDb, err := db.NewMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	db.Client = mongoDb

	// Initialize default config
	app.Use(logger.New(logger.Config{}))

	// Initialize router
	router.Init(&r)

	defer app.Server().Shutdown()

	if err := app.Listen(":3008"); err != nil {
		log.Fatal(err)
	}
}
