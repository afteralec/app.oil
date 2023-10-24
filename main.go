package main

import (
	"embed"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/queries"
)

//go:embed web/views/*
var viewsfs embed.FS

func main() {
	queries.Connect()
	queries.Build()
	config := configs.Fiber(viewsfs)
	app := fiber.New(config)

	app.Use(cors.New())
	app.Use(logger.New())

	app.Static("/", "./web/static")

	app.Get("/", handlers.Home)

	app.Get("/request/:id", handlers.Request)

	app.Listen(":8008")
}
