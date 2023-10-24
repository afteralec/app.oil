package main

import (
	"embed"
  "log"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
  
  "petrichormud.com/app/internal/middleware"
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
  app.Use(middleware.Session())

	app.Static("/", "./web/static")

	app.Get("/", handlers.Home)

  player := app.Group("player")
  player.Post("/", handlers.NewPlayer)

  request := app.Group("request")
	request.Get("/:id", handlers.Request)

	log.Fatal(app.Listen(":8008"))
}
