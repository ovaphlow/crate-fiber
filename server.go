package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Serve(addr string) {
	// init database

	app := fiber.New(fiber.Config{
		Prefork:   true,
		BodyLimit: 16 * 1024 * 1024,
	})

	app.Use(compress.New())

	app.Use(etag.New())

	app.Use(helmet.New())

	app.Use(logger.New(logger.Config{
		Format:     "${ip} [${time}] ${status} - ${latency} ${method} ${path} ?${queryParams}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Shanghai",
	}))

	app.Use(recover.New())

	app.Use(func(c *fiber.Ctx) error {
		// uri 过滤
		return c.Next()
	})

	app.Get("/crate-api/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "hola el mondo",
		})
	})

	app.Get("/crate-api/events", EventsEndpointGet)

	log.Fatal(app.Listen(addr))
}
